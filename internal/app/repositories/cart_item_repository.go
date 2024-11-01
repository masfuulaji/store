package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

type CartItemRepository interface {
	CreateCartItem(cartItem models.CartItem) error
	GetCartItem(id string) (models.CartItem, error)
	GetCartItems() ([]models.CartItem, error)
	GetCartItemsByCart(id string) ([]models.CartItem, error)
	GetCartItemByCart(id string) ([]models.CartItem, error)
	SumCartItemByCart(id string) (int, error)
	UpdateCartItem(cartItem models.CartItem, id string) error
	DeleteCartItem(id string) error
}

type CartItemRepositoryImpl struct {
	db *sqlx.DB
}

func NewCartItemRepository(db *sqlx.DB) *CartItemRepositoryImpl {
	return &CartItemRepositoryImpl{db: db}
}

func (u *CartItemRepositoryImpl) CreateCartItem(cartItem models.CartItem) error {
	query := "INSERT INTO cart_items (cart_id, product_id, product_qty, price_total, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, cartItem.CartId, cartItem.ProductId, cartItem.ProductQty, cartItem.PriceTotal, createdAt, updatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *CartItemRepositoryImpl) GetCartItem(id string) (models.CartItem, error) {
	var category models.CartItem
	query := "SELECT * FROM cart_items WHERE id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&category, query, id)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (u *CartItemRepositoryImpl) GetCartItems() ([]models.CartItem, error) {
	var cartItems []models.CartItem
	query := "SELECT * FROM cart_items WHERE deleted_at IS NULL"
	err := u.db.Select(&cartItems, query)
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

func (u *CartItemRepositoryImpl) GetCartItemsByCart(id string) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	query := "SELECT * FROM cart_items WHERE cart_id = $1 AND deleted_at IS NULL"
	err := u.db.Select(&cartItems, query, id)
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

func (u *CartItemRepositoryImpl) GetCartItemByCart(id string) (models.CartItem, error) {
	var cartItems models.CartItem
	query := "SELECT * FROM cart_items WHERE cart_id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&cartItems, query, id)
	if err != nil {
		return cartItems, err
	}
	return cartItems, nil
}
func (u *CartItemRepositoryImpl) SumCartItemByCart(id string) (int, error) {
	var count float64
	query := "SELECT SUM(price_total) FROM cart_items WHERE cart_id = $1 AND deleted_at IS NULL"
	err := u.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return int(count), err
	}
	return int(count), nil
}

func (u *CartItemRepositoryImpl) UpdateCartItem(cartItem models.CartItem, id string) error {
	query := "UPDATE cart_items SET product_qty = $1, price_total = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, cartItem.ProductQty, cartItem.PriceTotal, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *CartItemRepositoryImpl) DeleteCartItem(id string) error {
	deletedAt := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE cart_items SET deleted_at = $1 WHERE id = $2"
	_, err := u.db.Exec(query, deletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
