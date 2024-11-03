package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

type ProductRepository interface {
	CreateProduct(product models.Product) error
	GetProduct(id string) (models.Product, error)
	GetProducts() ([]models.Product, error)
	UpdateProduct(product models.Product, id string) error
	UpdateProductStock(stock int, id string) error
	DeleteProduct(id string) error
}

type ProductRepositoryImpl struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (u *ProductRepositoryImpl) CreateProduct(product models.Product) error {
	query := "INSERT INTO products (name, category_id, price, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, product.Name, product.CategoryId, product.Price, product.Stock, createdAt, updatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *ProductRepositoryImpl) GetProduct(id string) (models.Product, error) {
	var product models.Product
	query := "SELECT * FROM products WHERE id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&product, query, id)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (u *ProductRepositoryImpl) GetProducts() ([]models.Product, error) {
	var products []models.Product
	query := "SELECT * FROM products WHERE deleted_at IS NULL"
	err := u.db.Select(&products, query)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (u *ProductRepositoryImpl) GetProductsByCategory(id string) ([]models.Product, error) {
	var products []models.Product
	query := "SELECT * FROM products WHERE category_id = $1 AND deleted_at IS NULL"
	err := u.db.Select(&products, query, id)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (u *ProductRepositoryImpl) UpdateProduct(product models.Product, id string) error {
	query := "UPDATE products SET name = $1, category_id = $2, price = $3, stock = $4, updated_at = $5 WHERE id = $6"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, product.Name, product.CategoryId, product.Price, product.Stock, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *ProductRepositoryImpl) UpdateProductStock(stock int, id string) error {
	query := "UPDATE products SET stock = $1, updated_at = $2 WHERE id = $3"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, stock, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *ProductRepositoryImpl) DeleteProduct(id string) error {
	deletedAt := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE products SET deleted_at = $1 WHERE id = $2"
	_, err := u.db.Exec(query, deletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
