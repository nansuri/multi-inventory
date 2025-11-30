package postgres

import "time"

type UserModel struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string    `gorm:"type:text;uniqueIndex;not null"`
	Password  string    `gorm:"type:text;not null"`
	Role      string    `gorm:"type:text;not null;default:user"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

func (UserModel) TableName() string { return "users" }

type ItemModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:text;not null"`
	Barcode   string    `gorm:"type:text;uniqueIndex;not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	Location  string    `gorm:"type:text"`
	IsHalal   bool      `gorm:"not null;default:true"`
	Quantity  int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

func (ItemModel) TableName() string { return "items" }

type SalesOrderModel struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     *string   `gorm:"type:uuid"`
	TotalPrice float64   `gorm:"type:decimal(10,2);not null"`
	Status     string    `gorm:"type:text;not null;default:pending"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	UpdatedAt  time.Time `gorm:"not null;default:now()"`
}

func (SalesOrderModel) TableName() string { return "sales_orders" }

type SalesOrderItemModel struct {
	ID           int64   `gorm:"primaryKey;autoIncrement"`
	SalesOrderID int64   `gorm:"not null"`
	ItemID       int64   `gorm:"not null"`
	Quantity     int     `gorm:"not null"`
	PriceAtSale  float64 `gorm:"type:decimal(10,2);not null"`
	IsFulfilled  bool    `gorm:"not null;default:false"`
}

func (SalesOrderItemModel) TableName() string { return "sales_order_items" }
