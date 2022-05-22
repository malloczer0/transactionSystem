package pg

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"transactionSystemTestTask/model"
)

// TransactionPgRepo ...
type TransactionPgRepo struct {
	db *DB
}

// NewTransactionRepo ...
func NewTransactionRepo(db *DB) *TransactionPgRepo {
	return &TransactionPgRepo{db: db}
}

// GetTransaction retrieves transaction from Postgres
func (repo *TransactionPgRepo) GetTransaction(ctx context.Context, id uuid.UUID) (*model.DBTransaction, error) {
	transaction := &model.DBTransaction{}
	err := repo.db.Model(transaction).
		Where("id = ?", id).
		Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return transaction, nil
}

// GetPending retrieves all transactions with pending status from Postgres
func (repo *TransactionPgRepo) GetPending(ctx context.Context) (*[]model.DBTransaction, error) {
	transactions := &[]model.DBTransaction{}
	err := repo.db.Model(transactions).
		Where("status = ?", 0).
		Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return transactions, nil
}

// CreateTransaction creates transaction in Postgres
func (repo *TransactionPgRepo) CreateTransaction(ctx context.Context, transaction *model.DBTransaction) (*model.DBTransaction, error) {
	_, err := repo.db.Model(transaction).
		Returning("*").
		Insert()
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// UpdateTransaction updates transaction in Postgres
func (repo *TransactionPgRepo) UpdateTransaction(ctx context.Context, transaction *model.DBTransaction) (*model.DBTransaction, error) {
	_, err := repo.db.Model(transaction).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}

	return transaction, nil
}
