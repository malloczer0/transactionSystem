package model

import (
	"github.com/google/uuid"
	"time"
)

// Transaction is a JSON transaction
type Transaction struct {
	ID        uuid.UUID `json:"id"`
	ClientId  uuid.UUID `json:"client_id" validate:"required"`
	Status    int8      `json:"status" validate:"required:gte=-1,lte=1"`
	Change    float64   `json:"change"`
	CreatedAt time.Time `json:"created_at"`
}

// ToDB converts Transaction to DBTransaction
func (transaction *Transaction) ToDB() *DBTransaction {
	return &DBTransaction{
		ID:        transaction.ID,
		ClientId:  transaction.ClientId,
		Status:    transaction.Status,
		Change:    transaction.Change,
		CreatedAt: transaction.CreatedAt,
	}
}

// DBTransaction is a Postgres transaction
type DBTransaction struct {
	tableName struct{}  `pg:"transactions" gorm:"primaryKey"`
	ID        uuid.UUID `pg:"id,notnull,pk"`
	ClientId  uuid.UUID `pg:"client_id,notnull"`
	Status    int8      `pg:"status,notnull"`
	Change    float64   `pg:"change,notnull"`
	CreatedAt time.Time `pg:"created_at,notnull"`
}

// TableName overrides default table name for gorm
func (DBTransaction) TableName() string {
	return "transactions"
}

// ToWeb converts DBTransaction to Transaction
func (dbTransaction *DBTransaction) ToWeb() *Transaction {
	return &Transaction{
		ID:        dbTransaction.ID,
		ClientId:  dbTransaction.ClientId,
		Status:    dbTransaction.Status,
		Change:    dbTransaction.Change,
		CreatedAt: dbTransaction.CreatedAt,
	}
}
