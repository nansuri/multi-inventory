package postgres

import (
	"context"
	"errors"
	"fmt"
	"multi-inventory/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ItemRepository struct {
	db *DB
}

func NewItemRepository(db *DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *domain.Item) error {
	query := `
		INSERT INTO items (name, barcode, price, location, is_halal, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, item.Name, item.Barcode, item.Price, item.Location, item.IsHalal, item.Quantity).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			return fmt.Errorf("barcode already exists")
		}
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (r *ItemRepository) Update(ctx context.Context, item *domain.Item) error {
	query := `
		UPDATE items
		SET name = $1, barcode = $2, price = $3, location = $4, is_halal = $5, quantity = $6, updated_at = NOW()
		WHERE id = $7
	`
	cmdTag, err := r.db.Pool.Exec(ctx, query, item.Name, item.Barcode, item.Price, item.Location, item.IsHalal, item.Quantity, item.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			return fmt.Errorf("barcode already exists")
		}
		return fmt.Errorf("failed to update item: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = $1`
	cmdTag, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("item not found")
	}
	return nil
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	query := `
		SELECT id, name, barcode, price, location, is_halal, quantity, created_at, updated_at
		FROM items
		WHERE id = $1
	`
	var item domain.Item
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.Name, &item.Barcode, &item.Price, &item.Location, &item.IsHalal, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get item by id: %w", err)
	}
	return &item, nil
}

func (r *ItemRepository) GetByBarcode(ctx context.Context, barcode string) (*domain.Item, error) {
	query := `
		SELECT id, name, barcode, price, location, is_halal, quantity, created_at, updated_at
		FROM items
		WHERE barcode = $1
	`
	var item domain.Item
	err := r.db.Pool.QueryRow(ctx, query, barcode).Scan(
		&item.ID, &item.Name, &item.Barcode, &item.Price, &item.Location, &item.IsHalal, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get item by barcode: %w", err)
	}
	return &item, nil
}

func (r *ItemRepository) List(ctx context.Context) ([]*domain.Item, error) {
	query := `
		SELECT id, name, barcode, price, location, is_halal, quantity, created_at, updated_at
		FROM items
		ORDER BY name ASC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Barcode, &item.Price, &item.Location, &item.IsHalal, &item.Quantity, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}
	return items, nil
}
