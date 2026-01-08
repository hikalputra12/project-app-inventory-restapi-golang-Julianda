package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type TransactionRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type TransactionRepoInterface interface {
	GetAllTransaction(page, limit int) ([]model.Transaction, int, error)
	CreateTransaction(Transaction *model.Transaction) error
	UpdateTransaction(id int, Transaction *model.Transaction) error
	DeleteTransaction(id int) error
}

// constructor
func NewTransactionRepo(db database.PgxIface,
	log *zap.Logger) TransactionRepoInterface {
	return &TransactionRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *TransactionRepo) CreateTransaction(transaction *model.Transaction) error {
	// transaksi untuk insert dan edit secara bersamaan dengan menggunakan tx
	tx, err := r.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	queryInsert := `
        INSERT INTO sales_item (user_id, inventory_id, quantity, price, created_at, updated_at)
SELECT $1, $2, $3, price, $4, $5  
FROM inventories 
WHERE inventory_id = $2        
RETURNING sales_item_id
    `
	now := time.Now()

	err = tx.QueryRow(context.Background(), queryInsert,
		transaction.UserId,
		transaction.InventoryId,
		transaction.Quantity,
		now, now,
	).Scan(&transaction.ID)

	if err != nil {
		r.Logger.Error("Database Query Error: Gagal create transaksi penjualan item inventory",
			zap.Error(err),
			zap.String("query", queryInsert),
		)
		return err
	}

	//query untuk edit stock di inventory
	queryUpdate := `
        UPDATE inventories 
        SET stock = stock - $1, updated_at = $2
        WHERE inventory_id = $3 AND stock >= $1
    `

	cmdTag, err := tx.Exec(context.Background(), queryUpdate, transaction.Quantity, now, transaction.InventoryId)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal mengubah stok item di tabel inventoris",
			zap.Error(err),
			zap.String("query", queryUpdate),
		)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("gagal update stok: stok tidak cukup atau barang tidak ditemukan")
	}

	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return tx.Commit(context.Background())
}

// untuk membaca Transaction yang ada
func (r *TransactionRepo) GetAllTransaction(page, limit int) ([]model.Transaction, int, error) {

	offset := (page - 1) * limit
	var total int
	countQuery := `SELECT COUNT(*) FROM sales_item WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT i.name,t.quantity,t.price FROM sales_item t
	JOIN inventories i ON t.inventory_id = i.inventory_id
	ORDER BY sales_item_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	var transaction []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.Name, &t.Quantity, &t.Price)
		if err != nil {
			return nil, 0, err
		}
		transaction = append(transaction, t)
	}
	return transaction, total, nil
}

func (r *TransactionRepo) UpdateTransaction(id int, transaction *model.Transaction) error {
	ctx := context.Background()

	// 1. MULAI TRANSAKSI
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// AMBIL DATA LAMA
	// Kita perlu tahu inventory_id mana yang diedit & berapa jumlah awalnya
	var oldQty int
	var inventoryID int

	// Gunakan FOR UPDATE untuk mengunci baris agar aman dari race condition
	queryGetOld := `SELECT quantity, inventory_id FROM sales_item WHERE sales_item_id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, queryGetOld, id).Scan(&oldQty, &inventoryID)
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan: %w", err)
	}

	// HITUNG SELISIH (PERUBAHAN)
	// Rumus: Baru - Lama
	change := transaction.Quantity - oldQty

	// Jika tidak ada perubahan jumlah, langsung return (hemat proses DB)
	if change == 0 {
		return nil
	}

	now := time.Now()

	// UPDATE SALES_ITEM (Simpan jumlah baru)
	queryUpdateItem := `UPDATE sales_item SET quantity=$1, updated_at=$2 WHERE sales_item_id=$3`
	_, err = tx.Exec(ctx, queryUpdateItem, transaction.Quantity, now, id)
	if err != nil {
		return err
	}

	// UPDATE STOCK INVENTORY (Gunakan Selisih)
	// Rumus: stock = stock - change
	// Jika change positif (nambah beli 2): stock - 2 (Stok berkurang)
	// Jika change negatif (batal beli 2): stock - (-2) => stock + 2 (Stok balik)
	queryUpdateStock := `
        UPDATE inventories 
        SET stock = stock - $1, updated_at = $2
        WHERE inventory_id = $3 AND stock >= $1
    `

	// Perhatikan parameter: $1 diisi 'change' (selisih), bukan total quantity
	cmdTag, err := tx.Exec(ctx, queryUpdateStock, change, now, inventoryID)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal update stock transaksi",
			zap.Error(err),
			zap.String("query", queryUpdateStock),
		)
		return err
	}

	// Validasi: Jika change positif (mengurangi stok) tapi stok tidak cukup
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("gagal update: stok gudang tidak mencukupi untuk penambahan jumlah")
	}

	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	// 6. COMMIT
	return tx.Commit(ctx)
}

// delete Transaction
func (r *TransactionRepo) DeleteTransaction(id int) error {
	query := `DELETE FROM sales_item
			 where sales_item_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal cmebnghapus transaksi",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}
