package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

type CategoryRepository interface {
	CreateCategory(category models.Category) (int, error)
	GetCategory(id string) (models.Category, error)
	GetCategories() ([]models.Category, error)
	UpdateCategory(category models.Category, id string) error
	DeleteCategory(id string) error
}

type CategoryRepositoryImpl struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db}
}

func (u *CategoryRepositoryImpl) CreateCategory(category models.Category) (int, error) {
	query := "INSERT INTO categories (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err := u.db.QueryRow(query, category.Name, createdAt, updatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *CategoryRepositoryImpl) GetCategory(id string) (models.Category, error) {
	var category models.Category
	query := "SELECT * FROM categories WHERE id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&category, query, id)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (u *CategoryRepositoryImpl) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	query := "SELECT * FROM categories WHERE deleted_at IS NULL"
	err := u.db.Select(&categories, query)
	if err != nil {
		return categories, err
	}
	return categories, nil
}

func (u *CategoryRepositoryImpl) UpdateCategory(category models.Category, id string) error {
	query := "UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, category.Name, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *CategoryRepositoryImpl) DeleteCategory(id string) error {
	deletedAt := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE categories SET deleted_at = $1 WHERE id = $2"
	_, err := u.db.Exec(query, deletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
