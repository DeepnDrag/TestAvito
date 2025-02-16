package storage

import (
	"TestAvito/internal/models"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (s *ProductRepo) GetItemPrice(productName string) (int, error) {
	var product models.Product

	err := s.db.Where("name = ?", productName).First(&product).Error
	if err != nil {
		return 0, err
	}

	return product.Price, nil
}
