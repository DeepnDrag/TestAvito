package storage

import (
	"TestAvito/internal/models"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (s *TransactionRepo) CreateTransaction(fromUserID, toUserID uint, amount int) (*models.Transaction, error) {
	transaction := models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}
	if err := s.db.Create(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *TransactionRepo) GetGiftsGivenByUser(userID uint) ([]models.TransactionsFromUser, error) {
	var result []models.TransactionsFromUser

	err := s.db.Table("transactions").
		Select("to_user_id AS to_user, SUM(amount) AS amount").
		Where("from_user_id = ?", userID).
		Group("to_user_id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *TransactionRepo) GetGiftsGivenToUser(userID uint) ([]models.TransactionsToUser, error) {
	var result []models.TransactionsToUser

	err := s.db.Table("transactions").
		Select("from_user_id AS from_user, SUM(amount) AS amount").
		Where("to_user_id = ?", userID).
		Group("from_user_id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, err
}
