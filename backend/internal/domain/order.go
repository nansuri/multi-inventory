package domain

import (
	"context"
	"time"
)

type SalesOrder struct {
	ID         int64             `json:"id"`
	UserID     int64             `json:"user_id"`
	TotalPrice float64           `json:"total_price"`
	Status     string            `json:"status"` // pending, completed, cancelled
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Items      []*SalesOrderItem `json:"items,omitempty"`
}

type SalesOrderItem struct {
	ID           int64   `json:"id"`
	SalesOrderID int64   `json:"sales_order_id"`
	ItemID       int64   `json:"item_id"`
	ItemName     string  `json:"item_name,omitempty"` // For display
	Quantity     int     `json:"quantity"`
	PriceAtSale  float64 `json:"price_at_sale"`
	IsFulfilled  bool    `json:"is_fulfilled"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *SalesOrder) error
	GetByID(ctx context.Context, id int64) (*SalesOrder, error)
	List(ctx context.Context) ([]*SalesOrder, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateItemFulfillment(ctx context.Context, itemId int64, isFulfilled bool) error
}
