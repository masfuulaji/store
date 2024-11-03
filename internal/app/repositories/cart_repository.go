package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

type CartRepository interface {
	CreateCart(cart models.Cart) (int, error)
	GetCart(id string) (models.Cart, error)
	GetCartByUserId(id string) (models.Cart, error)
	CountCartByUserId(id string) (int, error)
	GetCarts() ([]models.Cart, error)
	UpdateCart(cart models.Cart, id string) error
	UpdateCartTotal(sum int, id string) error
	DeleteCart(id string) error
}

type CartRepositoryImpl struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepositoryImpl {
	return &CartRepositoryImpl{db: db}
}

func (u *CartRepositoryImpl) CreateCart(cart models.Cart) (int, error) {
	query := "INSERT INTO carts (name, user_id, price_total, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err := u.db.QueryRow(query, cart.Name, cart.UserId, cart.PriceTotal, createdAt, updatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *CartRepositoryImpl) GetCart(id string) (models.Cart, error) {
	var cart models.Cart
	query := "SELECT * FROM carts WHERE id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&cart, query, id)
	if err != nil {
		return cart, err
	}
	return cart, nil
}

func (u *CartRepositoryImpl) GetCartByUserId(id string) (models.Cart, error) {
	var cart models.Cart
	query := "SELECT * FROM carts WHERE user_id = $1 AND finish IS NULL AND deleted_at IS NULL"
	err := u.db.Get(&cart, query, id)
	if err != nil {
		return cart, err
	}
	return cart, nil
}

func (u *CartRepositoryImpl) CountCartByUserId(id string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM carts WHERE user_id = $1 AND finish is null AND deleted_at IS NULL"
	err := u.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (u *CartRepositoryImpl) GetCarts() ([]models.Cart, error) {
	var carts []models.Cart
	query := "SELECT * FROM carts WHERE deleted_at IS NULL"
	err := u.db.Select(&carts, query)
	if err != nil {
		return carts, err
	}
	return carts, nil
}

func (u *CartRepositoryImpl) UpdateCart(cart models.Cart, id string) error {
	query := "UPDATE carts SET name = $1, price_total = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, cart.Name, cart.PriceTotal, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *CartRepositoryImpl) UpdateCartTotal(sum int, id string) error {
	query := "UPDATE carts SET price_total = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, sum, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *CartRepositoryImpl) DeleteCart(id string) error {
	deletedAt := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE carts SET deleted_at = $1 WHERE id = $2"
	_, err := u.db.Exec(query, deletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
