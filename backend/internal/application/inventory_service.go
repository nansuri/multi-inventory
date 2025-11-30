package application

import (
	"context"
	"multi-inventory/internal/domain"
)

type InventoryService struct {
	itemRepo domain.ItemRepository
}

func NewInventoryService(itemRepo domain.ItemRepository) *InventoryService {
	return &InventoryService{itemRepo: itemRepo}
}

func (s *InventoryService) CreateItem(ctx context.Context, item *domain.Item) error {
	// Add validation logic here if needed
	return s.itemRepo.Create(ctx, item)
}

func (s *InventoryService) UpdateItem(ctx context.Context, item *domain.Item) error {
	return s.itemRepo.Update(ctx, item)
}

func (s *InventoryService) DeleteItem(ctx context.Context, id int64) error {
	return s.itemRepo.Delete(ctx, id)
}

func (s *InventoryService) GetItem(ctx context.Context, id int64) (*domain.Item, error) {
	return s.itemRepo.GetByID(ctx, id)
}

func (s *InventoryService) GetItemByBarcode(ctx context.Context, barcode string) (*domain.Item, error) {
	return s.itemRepo.GetByBarcode(ctx, barcode)
}

func (s *InventoryService) ListItems(ctx context.Context) ([]*domain.Item, error) {
	return s.itemRepo.List(ctx)
}
