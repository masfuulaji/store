package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

type OrderRepository interface {
	CreateOrder(order models.Order) (int, error)
}

type OrderRepositoryImpl struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

func (u *OrderRepositoryImpl) CreateOrder(order models.Order) (int, error) {
	query := "INSERT INTO orders (cart_id, price_total, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err := u.db.QueryRow(query, order.CartId, order.PriceTotal, createdAt, updatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
