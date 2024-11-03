package models

import (
	"time"

	"github.com/lib/pq"
)

type Order struct {
	ID         string      `json:"id"`
	CartId     string      `db:"cart_id" json:"cart_id"`
	PriceTotal float64     `db:"price_total" json:"price_total"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt  pq.NullTime `db:"deleted_at" json:"deleted_at"`
}
