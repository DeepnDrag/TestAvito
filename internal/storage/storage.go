package storage

import (
	"TestAvito/internal/models"
	"gorm.io/gorm"
)

type UserStorage interface {
	CreateUser(username, password string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(updatedUser *models.User) (*models.User, error)
	UpdateTwoUsers(updatedUser1 *models.User, updatedUser2 *models.User) (*models.User, *models.User, error)
}

type TransactionStorage interface {
	CreateTransaction(fromUserID, toUserID uint, amount int) (*models.Transaction, error)
	GetGiftsGivenByUser(userID uint) ([]models.TransactionsFromUser, error)
	GetGiftsGivenToUser(userID uint) ([]models.TransactionsToUser, error)
}

type InventoryStorage interface {
	CreateInventory(userID uint, itemType string, quantity int) (*models.Inventory, error)
	UpdateInventory(userID uint, itemType string, quantity int) (*models.Inventory, error)
	GetPurchasedItems(userID uint) ([]models.Inventory, error)
}

type ProductStorage interface {
	GetItemPrice(productName string) (int, error)
}

type Storage struct {
	UserStorage
	TransactionStorage
	InventoryStorage
	ProductStorage
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		UserStorage:        NewUserRepo(db),
		TransactionStorage: NewTransactionRepo(db),
		InventoryStorage:   NewInventoryRepo(db),
		ProductStorage:     NewProductRepo(db),
	}
}
