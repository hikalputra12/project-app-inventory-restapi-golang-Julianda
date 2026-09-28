package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type TransactionRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type TransactionRepoInterface interface {
	GetAllTransaction(ctx context.Context, page, limit int) ([]model.Transaction, int, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	UpdateTransaction(ctx context.Context, id int, transaction *model.Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
	GetTransactionById(ctx context.Context, id int) (*model.Transaction, error)
}

func NewTransactionRepo(db database.PgxIface, log *zap.Logger) TransactionRepoInterface {
	return &TransactionRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Lock inventory row and verify current stock
	var currentStock int
	var currentPrice int
	queryCheck := `SELECT stock, price FROM inventories WHERE inventory_id = $1 AND deleted_at IS NULL FOR UPDATE`
	err = tx.QueryRow(ctx, queryCheck, transaction.InventoryId).Scan(&currentStock, &currentPrice)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("barang inventaris dengan ID %d tidak ditemukan", transaction.InventoryId)
		}
		return fmt.Errorf("gagal mengecek inventaris: %w", err)
	}

	if currentStock < transaction.Quantity {
		return fmt.Errorf("stok barang tidak mencukupi (tersedia: %d, diminta: %d)", currentStock, transaction.Quantity)
	}

	now := time.Now()
	// 2. Insert into sales_item
	queryInsert := `
		INSERT INTO sales_item (user_id, inventory_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING sales_item_id
	`
	err = tx.QueryRow(ctx, queryInsert,
		transaction.UserId,
		transaction.InventoryId,
		transaction.Quantity,
		currentPrice,
		now, now,
	).Scan(&transaction.ID)
	if err != nil {
		r.Logger.Error("Failed to insert sales item", zap.Error(err))
		return err
	}

	// 3. Deduct stock safely
	queryUpdateStock := `
		UPDATE inventories 
		SET stock = stock - $1, updated_at = $2
		WHERE inventory_id = $3
	`
	_, err = tx.Exec(ctx, queryUpdateStock, transaction.Quantity, now, transaction.InventoryId)
	if err != nil {
		r.Logger.Error("Failed to update inventory stock", zap.Error(err))
		return err
	}

	transaction.Price = currentPrice
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return tx.Commit(ctx)
}

func (r *TransactionRepo) GetAllTransaction(ctx context.Context, page, limit int) ([]model.Transaction, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM sales_item WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT 
		t.sales_item_id,
		COALESCE(t.user_id, 0),
		t.inventory_id,
		i.name,
		t.quantity,
		t.price,
		t.created_at
	FROM sales_item t
	JOIN inventories i ON t.inventory_id = i.inventory_id
	WHERE t.deleted_at IS NULL
	ORDER BY t.sales_item_id DESC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.UserId, &t.InventoryId, &t.Name, &t.Quantity, &t.Price, &t.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepo) UpdateTransaction(ctx context.Context, id int, transaction *model.Transaction) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Lock existing sales item
	var oldQty int
	var inventoryID int
	queryGetOld := `SELECT quantity, inventory_id FROM sales_item WHERE sales_item_id = $1 AND deleted_at IS NULL FOR UPDATE`
	err = tx.QueryRow(ctx, queryGetOld, id).Scan(&oldQty, &inventoryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("transaksi tidak ditemukan atau sudah dihapus")
		}
		return fmt.Errorf("gagal memeriksa transaksi: %w", err)
	}

	change := transaction.Quantity - oldQty
	if change == 0 {
		return nil
	}

	// 2. Lock inventory row
	var currentStock int
	queryCheckStock := `SELECT stock FROM inventories WHERE inventory_id = $1 AND deleted_at IS NULL FOR UPDATE`
	err = tx.QueryRow(ctx, queryCheckStock, inventoryID).Scan(&currentStock)
	if err != nil {
		return fmt.Errorf("gagal memeriksa stok barang: %w", err)
	}

	// If increasing purchase quantity, check stock
	if change > 0 && currentStock < change {
		return fmt.Errorf("stok gudang tidak mencukupi untuk penambahan jumlah (stok tersisa: %d)", currentStock)
	}

	now := time.Now()

	// 3. Update sales_item quantity
	queryUpdateItem := `UPDATE sales_item SET quantity = $1, updated_at = $2 WHERE sales_item_id = $3`
	_, err = tx.Exec(ctx, queryUpdateItem, transaction.Quantity, now, id)
	if err != nil {
		return err
	}

	// 4. Update inventory stock
	queryUpdateStock := `UPDATE inventories SET stock = stock - $1, updated_at = $2 WHERE inventory_id = $3`
	_, err = tx.Exec(ctx, queryUpdateStock, change, now, inventoryID)
	if err != nil {
		return err
	}

	transaction.InventoryId = inventoryID
	transaction.UpdatedAt = now

	return tx.Commit(ctx)
}

func (r *TransactionRepo) DeleteTransaction(ctx context.Context, id int) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Lock and retrieve transaction details
	var qty int
	var inventoryID int
	queryGet := `SELECT quantity, inventory_id FROM sales_item WHERE sales_item_id = $1 AND deleted_at IS NULL FOR UPDATE`
	err = tx.QueryRow(ctx, queryGet, id).Scan(&qty, &inventoryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("transaksi tidak ditemukan atau sudah dihapus")
		}
		return err
	}

	now := time.Now()

	// 2. Soft delete sales_item
	querySoftDelete := `UPDATE sales_item SET deleted_at = $1, updated_at = $1 WHERE sales_item_id = $2`
	_, err = tx.Exec(ctx, querySoftDelete, now, id)
	if err != nil {
		return err
	}

	// 3. RESTORE INVENTORY STOCK (Fixing critical stock loss bug!)
	queryRestoreStock := `UPDATE inventories SET stock = stock + $1, updated_at = $2 WHERE inventory_id = $3`
	_, err = tx.Exec(ctx, queryRestoreStock, qty, now, inventoryID)
	if err != nil {
		return fmt.Errorf("gagal mengembalikan stok inventaris: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *TransactionRepo) GetTransactionById(ctx context.Context, id int) (*model.Transaction, error) {
	var transaction model.Transaction
	query := `SELECT 
		s.sales_item_id,
		COALESCE(s.user_id, 0),
		s.inventory_id,
		i.name, 
		s.quantity, 
		s.price,
		s.created_at
	FROM sales_item s
	JOIN inventories i ON s.inventory_id = i.inventory_id
	WHERE s.sales_item_id = $1 AND s.deleted_at IS NULL;`

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.InventoryId,
		&transaction.Name,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
