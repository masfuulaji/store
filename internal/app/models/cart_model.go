package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Cart struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	UserId     string        `db:"user_id" json:"user_id"`
	PriceTotal float64       `db:"price_total" json:"price_total"`
	Finish     sql.NullInt32 `json:"finish"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt  pq.NullTime   `db:"deleted_at" json:"deleted_at"`
}
