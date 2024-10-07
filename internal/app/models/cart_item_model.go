package models

import (
	"time"

	"github.com/lib/pq"
)

type CartItem struct {
	ID         string      `json:"id"`
	CartId     string      `db:"cart_id" json:"cart_id"`
	ProductId  string      `db:"product_id" json:"product_id"`
	ProductQty string      `db:"product_qty" json:"product_qty"`
	PriceTotal float64     `db:"price_total" json:"price_total"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt  pq.NullTime `db:"deleted_at" json:"deleted_at"`
}
