package domain

import (
	"context"
	"time"
)

type Item struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Barcode   string    `json:"barcode"`
	Price     float64   `json:"price"`
	Location  string    `json:"location"`
	IsHalal   bool      `json:"is_halal"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemRepository interface {
	Create(ctx context.Context, item *Item) error
	Update(ctx context.Context, item *Item) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Item, error)
	GetByBarcode(ctx context.Context, barcode string) (*Item, error)
	List(ctx context.Context) ([]*Item, error)
}
