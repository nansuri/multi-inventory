package postgres

import (
	"context"
	"fmt"
	"multi-inventory/internal/domain"

	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.SalesOrder) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create Order
	salesOrdersTable := fmt.Sprintf("%s.sales_orders", r.db.Schema)
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`, salesOrdersTable)
	var argUser any
	if order.UserID == "" {
		argUser = nil // insert NULL if not provided
	} else {
		argUser = order.UserID
	}
	err = tx.QueryRow(ctx, query, argUser, order.TotalPrice, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Create Order Items
	salesOrderItemsTable := fmt.Sprintf("%s.sales_order_items", r.db.Schema)
	itemQuery := fmt.Sprintf(`
		INSERT INTO %s (sales_order_id, item_id, quantity, price_at_sale, is_fulfilled)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, salesOrderItemsTable)
	for _, item := range order.Items {
		err = tx.QueryRow(ctx, itemQuery, order.ID, item.ItemID, item.Quantity, item.PriceAtSale, item.IsFulfilled).Scan(&item.ID)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
		item.SalesOrderID = order.ID
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.SalesOrder, error) {
	salesOrdersTable := fmt.Sprintf("%s.sales_orders", r.db.Schema)
	query := fmt.Sprintf(`
		SELECT id, user_id::text, total_price, status, created_at, updated_at
		FROM %s
		WHERE id = $1
	`, salesOrdersTable)
	var order domain.SalesOrder
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get Items
	salesOrderItemsTable := fmt.Sprintf("%s.sales_order_items", r.db.Schema)
	itemsTable := fmt.Sprintf("%s.items", r.db.Schema)
	itemsQuery := fmt.Sprintf(`
		SELECT soi.id, soi.sales_order_id, soi.item_id, soi.quantity, soi.price_at_sale, soi.is_fulfilled, i.name
		FROM %s soi
		JOIN %s i ON soi.item_id = i.id
		WHERE soi.sales_order_id = $1
	`, salesOrderItemsTable, itemsTable)
	rows, err := r.db.Pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.SalesOrderItem
		if err := rows.Scan(&item.ID, &item.SalesOrderID, &item.ItemID, &item.Quantity, &item.PriceAtSale, &item.IsFulfilled, &item.ItemName); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		order.Items = append(order.Items, &item)
	}

	return &order, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]*domain.SalesOrder, error) {
	salesOrdersTable := fmt.Sprintf("%s.sales_orders", r.db.Schema)
	query := fmt.Sprintf(`
		SELECT id, user_id::text, total_price, status, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`, salesOrdersTable)
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.SalesOrder
	for rows.Next() {
		var order domain.SalesOrder
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE sales_orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *OrderRepository) UpdateItemFulfillment(ctx context.Context, itemId int64, isFulfilled bool) error {
	query := `UPDATE sales_order_items SET is_fulfilled = $1 WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, isFulfilled, itemId)
	return err
}
