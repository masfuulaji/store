package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/masfuulaji/store/config"
)

type DB struct {
	DB *sqlx.DB
}

func ConnectDB() (*DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("config error")
	}
	source := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name, cfg.Database.Password)
	db, err := sqlx.Connect("postgres", source)
	// db, err := sqlx.Connect("postgres", "host=db port=5432 user=postgres dbname=store password=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}
func (db *DB) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}
