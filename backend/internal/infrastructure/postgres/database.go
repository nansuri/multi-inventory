package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
)

type DB struct {
	Pool   *pgxpool.Pool
	Schema string
	Gorm   *gorm.DB
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

	dbSchema := os.Getenv("DB_SCHEMA")
	if dbSchema == "" {
		dbSchema = "public"
	}

	gormDB, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{
		NamingStrategy: gormschema.NamingStrategy{
			TablePrefix: dbSchema + ".",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to open gorm: %w", err)
	}

	return &DB{Pool: pool, Schema: dbSchema, Gorm: gormDB}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

// AutoMigrate creates core tables if they do not exist.
// This is a simple MVP-friendly migration without external tools.
func (db *DB) AutoMigrate(ctx context.Context) error {
	// Ensure schema and extension
	if _, err := db.Pool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", db.Schema)); err != nil {
		return fmt.Errorf("failed to ensure schema: %w", err)
	}
	if err := db.Gorm.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		return fmt.Errorf("failed to ensure pgcrypto: %w", err)
	}

	// GORM automigrate
	if err := db.Gorm.WithContext(ctx).AutoMigrate(
		&UserModel{},
		&ItemModel{},
		&SalesOrderModel{},
		&SalesOrderItemModel{},
	); err != nil {
		return fmt.Errorf("gorm automigrate failed: %w", err)
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
