package pg

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"transactionSystemTestTask/model"
)

// ClientPgRepo ...
type ClientPgRepo struct {
	db *DB
}

// NewTransactionRepo ...
func NewClientRepo(db *DB) *ClientPgRepo {
	return &ClientPgRepo{db: db}
}

// GetTransaction retrieves client from Postgres
func (repo *ClientPgRepo) GetClient(ctx context.Context, id uuid.UUID) (*model.DBClient, error) {
	client := &model.DBClient{}
	err := repo.db.Model(client).
		Where("id = ?", id).
		Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return client, nil
}

// CreateTransaction creates client in Postgres
func (repo *ClientPgRepo) CreateClient(ctx context.Context, client *model.DBClient) (*model.DBClient, error) {
	_, err := repo.db.Model(client).
		Returning("*").
		Insert()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// UpdateClient updates client in Postgres
func (repo *ClientPgRepo) UpdateClient(ctx context.Context, client *model.DBClient) (*model.DBClient, error) {
	_, err := repo.db.Model(client).
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}

	return client, nil
}

// DeleteTransaction deletes client in Postgres
func (repo *ClientPgRepo) DeleteClient(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Model((*model.DBClient)(nil)).
		Where("id = ?", id).
		Delete()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
