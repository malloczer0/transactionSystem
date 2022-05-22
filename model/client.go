package model

import (
	"github.com/google/uuid"
	"time"
)

// Client is a JSON user
type Client struct {
	ID        uuid.UUID `json:"id"`
	Bio       string    `json:"bio" validate:"required"`
	Balance   float64   `json:"balance" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// ToDB converts Client to DBClient
func (client *Client) ToDB() *DBClient {
	return &DBClient{
		ID:        client.ID,
		Bio:       client.Bio,
		Balance:   client.Balance,
		CreatedAt: client.CreatedAt,
	}
}

// DBClient is a Postgres client
type DBClient struct {
	tableName struct{}  `pg:"clients" gorm:"primaryKey"`
	ID        uuid.UUID `pg:"id,notnull,pk"`
	Bio       string    `pg:"bio,notnull"`
	Balance   float64   `pg:"balance,notnull"`
	CreatedAt time.Time `pg:"created_at,notnull"`
}

// TableName overrides default table name for gorm
func (DBClient) TableName() string {
	return "clients"
}

// ToWeb converts DBClient to Client
func (dbClient *DBClient) ToWeb() *Client {
	return &Client{
		ID:        dbClient.ID,
		Bio:       dbClient.Bio,
		Balance:   dbClient.Balance,
		CreatedAt: dbClient.CreatedAt,
	}
}
