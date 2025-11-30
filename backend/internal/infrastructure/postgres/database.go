package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() (*DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbPort == "" {
		dbPort = "5432"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

// AutoMigrate creates core tables if they do not exist.
// This is a simple MVP-friendly migration without external tools.
func (db *DB) AutoMigrate(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS items (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			barcode TEXT NOT NULL UNIQUE,
			price DECIMAL(10, 2) NOT NULL,
			location TEXT,
			is_halal BOOLEAN DEFAULT TRUE,
			quantity INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS sales_orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id),
			total_price DECIMAL(10, 2) NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS sales_order_items (
			id BIGSERIAL PRIMARY KEY,
			sales_order_id BIGINT REFERENCES sales_orders(id) ON DELETE CASCADE,
			item_id BIGINT REFERENCES items(id),
			quantity INTEGER NOT NULL,
			price_at_sale DECIMAL(10, 2) NOT NULL,
			is_fulfilled BOOLEAN DEFAULT FALSE
		);`,
	}

	for _, sql := range stmts {
		if _, err := db.Pool.Exec(ctx, sql); err != nil {
			return fmt.Errorf("automigrate failed: %w", err)
		}
	}
	return nil
}
