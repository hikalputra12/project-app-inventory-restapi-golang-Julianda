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
	CreateTransaction(transaction *model.Transaction) error
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
		return err // Otomatis Rollback karena defer di atas
	}

	//query untuk edit stock di inventory
	queryUpdate := `
        UPDATE inventories 
        SET stock = stock - $1, updated_at = $2
        WHERE inventory_id = $3 AND stock >= $1
    `

	cmdTag, err := tx.Exec(context.Background(), queryUpdate, transaction.Quantity, now, transaction.InventoryId)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("gagal update stok: stok tidak cukup atau barang tidak ditemukan")
	}

	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	return tx.Commit(context.Background())
}
