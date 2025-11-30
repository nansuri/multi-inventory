package application

import (
	"context"
	"fmt"
	"multi-inventory/internal/domain"
)

type SalesService struct {
	orderRepo domain.OrderRepository
	itemRepo  domain.ItemRepository
}

func NewSalesService(orderRepo domain.OrderRepository, itemRepo domain.ItemRepository) *SalesService {
	return &SalesService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

func (s *SalesService) CreateOrder(ctx context.Context, userID int64, items []struct {
	ItemID   int64 `json:"item_id"`
	Quantity int   `json:"quantity"`
}) (*domain.SalesOrder, error) {
	order := &domain.SalesOrder{
		UserID: userID,
		Status: "pending",
		Items:  make([]*domain.SalesOrderItem, 0, len(items)),
	}

	var totalPrice float64

	for _, reqItem := range items {
		// Fetch item to get price and check availability
		item, err := s.itemRepo.GetByID(ctx, reqItem.ItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get item %d: %w", reqItem.ItemID, err)
		}
		if item == nil {
			return nil, fmt.Errorf("item %d not found", reqItem.ItemID)
		}
		if item.Quantity < reqItem.Quantity {
			return nil, fmt.Errorf("insufficient quantity for item %s", item.Name)
		}

		// Update item quantity (simple approach, should be transactional ideally)
		item.Quantity -= reqItem.Quantity
		if err := s.itemRepo.Update(ctx, item); err != nil {
			return nil, fmt.Errorf("failed to update inventory for item %s: %w", item.Name, err)
		}

		orderItem := &domain.SalesOrderItem{
			ItemID:      item.ID,
			Quantity:    reqItem.Quantity,
			PriceAtSale: item.Price,
			IsFulfilled: false,
		}
		order.Items = append(order.Items, orderItem)
		totalPrice += item.Price * float64(reqItem.Quantity)
	}

	order.TotalPrice = totalPrice

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *SalesService) ListOrders(ctx context.Context) ([]*domain.SalesOrder, error) {
	return s.orderRepo.List(ctx)
}

func (s *SalesService) GetOrder(ctx context.Context, id int64) (*domain.SalesOrder, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *SalesService) UpdateItemFulfillment(ctx context.Context, itemId int64, isFulfilled bool) error {
	return s.orderRepo.UpdateItemFulfillment(ctx, itemId, isFulfilled)
}
