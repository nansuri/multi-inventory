package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool   *pgxpool.Pool
	Schema string
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

	schema := os.Getenv("DB_SCHEMA")
	if schema == "" {
		schema = "public"
	}

	return &DB{Pool: pool, Schema: schema}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

// AutoMigrate creates core tables if they do not exist.
// This is a simple MVP-friendly migration without external tools.
func (db *DB) AutoMigrate(ctx context.Context) error {
	// Detect existing ID types (to avoid FK dtype mismatch with pre-existing tables)
	userIDType := detectColumnType(ctx, db, "users", "id", "bigint")
	itemIDType := detectColumnType(ctx, db, "items", "id", "bigint")

	// Ensure schema exists
	if _, err := db.Pool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", db.Schema)); err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}

	// Create base tables first (users, items) in target schema
	baseTables := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.users (
			id BIGSERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`, db.Schema),
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.items (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			barcode TEXT NOT NULL UNIQUE,
			price DECIMAL(10, 2) NOT NULL,
			location TEXT,
			is_halal BOOLEAN DEFAULT TRUE,
			quantity INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`, db.Schema),
	}
	for _, sql := range baseTables {
		if _, err := db.Pool.Exec(ctx, sql); err != nil {
			return fmt.Errorf("automigrate base tables failed: %w", err)
		}
	}

	// Normalize detected types for FK columns (Postgres uses 'integer' or 'bigint')
	fkUserType := mapTypeForFK(userIDType)
	fkItemType := mapTypeForFK(itemIDType)

	salesOrdersSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.sales_orders (
		id BIGSERIAL PRIMARY KEY,
		user_id %s REFERENCES %s.users(id),
		total_price DECIMAL(10, 2) NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`, db.Schema, fkUserType, db.Schema)

	salesOrderItemsSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.sales_order_items (
		id BIGSERIAL PRIMARY KEY,
		sales_order_id BIGINT REFERENCES %s.sales_orders(id) ON DELETE CASCADE,
		item_id %s REFERENCES %s.items(id),
		quantity INTEGER NOT NULL,
		price_at_sale DECIMAL(10, 2) NOT NULL,
		is_fulfilled BOOLEAN DEFAULT FALSE
	);`, db.Schema, db.Schema, fkItemType, db.Schema)

	if _, err := db.Pool.Exec(ctx, salesOrdersSQL); err != nil {
		return fmt.Errorf("automigrate sales_orders failed: %w (user id type detected: %s)", err, userIDType)
	}
	if _, err := db.Pool.Exec(ctx, salesOrderItemsSQL); err != nil {
		return fmt.Errorf("automigrate sales_order_items failed: %w (item id type detected: %s)", err, itemIDType)
	}
	return nil
}

// detectColumnType queries information_schema for an existing column data type.
func detectColumnType(ctx context.Context, db *DB, table, column, fallback string) string {
	var dtype string
	row := db.Pool.QueryRow(ctx, `SELECT data_type FROM information_schema.columns WHERE table_schema=$1 AND table_name=$2 AND column_name=$3`, db.Schema, table, column)
	if err := row.Scan(&dtype); err != nil || dtype == "" {
		return fallback
	}
	return dtype
}

// mapTypeForFK converts information_schema data_type to a column type usable in FK column definitions.
func mapTypeForFK(dtype string) string {
	switch dtype {
	case "integer":
		return "INTEGER"
	case "bigint":
		return "BIGINT"
	case "uuid":
		return "UUID"
	default:
		// Default to BIGINT for safety; Postgres will error if mismatched.
		return "BIGINT"
	}
}
