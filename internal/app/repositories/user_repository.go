package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masfuulaji/store/internal/app/models"
)

var user models.User

type UserRepository interface {
	CreateUser(user models.User) error
	GetUser(id string) (models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUser(user models.User, id string) error
	DeleteUser(id string) error
	GetUserByUsername(username string) (models.User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) CreateUser(user models.User) error {
	query := "INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, user.Username, user.Email, user.Password, createdAt, updatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryImpl) GetUser(id string) (models.User, error) {
	query := "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL"
	err := u.db.Get(&user, query, id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) GetUsers() ([]models.User, error) {
	query := "SELECT * FROM users WHERE deleted_at IS NULL"
	var users []models.User
	err := u.db.Select(&users, query)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (u *UserRepositoryImpl) GetUserByUsername(username string) (models.User, error) {
	query := "SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL"
	err := u.db.Get(&user, query, username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) UpdateUser(user models.User, id string) error {
	query := "UPDATE users SET username = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5"
	updatedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, user.Username, user.Email, user.Password, updatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepositoryImpl) DeleteUser(id string) error {
	query := "UPDATE users SET deleted_at = $1 WHERE id = $2"
	deletedAt := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.db.Exec(query, deletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
