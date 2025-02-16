package storage

import (
	"TestAvito/internal/models"
	"gorm.io/gorm"
)

type InventoryRepo struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) *InventoryRepo {
	return &InventoryRepo{
		db: db,
	}
}

func (s *InventoryRepo) CreateInventory(userID uint, itemType string, quantity int) (*models.Inventory, error) {
	newInventory := &models.Inventory{
		UserID:   userID,
		ItemType: itemType,
		Quantity: quantity,
	}

	err := s.db.Create(newInventory).Error
	if err != nil {
		return nil, err
	}

	return newInventory, nil
}

func (s *InventoryRepo) UpdateInventory(userID uint, itemType string, quantity int) (*models.Inventory, error) {
	var existingInventory models.Inventory
	err := s.db.Where("user_id = ? AND item_type = ?", userID, itemType).First(&existingInventory).Error
	if err != nil {
		return nil, err
	}

	existingInventory.Quantity += quantity
	err = s.db.Save(&existingInventory).Error
	if err != nil {
		return nil, err
	}

	return &existingInventory, nil
}

func (s *InventoryRepo) GetPurchasedItems(userID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory

	err := s.db.Where("user_id = ?", userID).Find(&inventory).Error
	if err != nil {
		return nil, err
	}

	return inventory, nil
}
