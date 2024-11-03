package models

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	CategoryId string      `db:"category_id" json:"category_id"`
	Price      float64     `json:"price"`
	Stock      int         `json:"stock"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt  pq.NullTime `db:"deleted_at" json:"deleted_at"`
}
