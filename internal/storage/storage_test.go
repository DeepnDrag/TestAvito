package storage

import (
	"TestAvito/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func getTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=testdb port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func applyMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Inventory{}, &models.Product{})
	if err != nil {
		slog.Info("Failed to apply migrations: %v", err)
	}
}

func clearTables(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users CASCADE")
	db.Exec("TRUNCATE TABLE transactions CASCADE")
	db.Exec("TRUNCATE TABLE inventories CASCADE")
	db.Exec("TRUNCATE TABLE products CASCADE")
}

func TestInventoryRepo_CreateInventory(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewInventoryRepo(db)
	inventory, err := repo.CreateInventory(1, "hoody", 5)
	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, uint(1), inventory.UserID)
	assert.Equal(t, "hoody", inventory.ItemType)
	assert.Equal(t, 5, inventory.Quantity)
}

func TestInventoryRepo_UpdateInventory(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewInventoryRepo(db)
	_, _ = repo.CreateInventory(1, "item1", 5)
	updatedInventory, err := repo.UpdateInventory(1, "item1", 3)
	assert.NoError(t, err)
	assert.NotNil(t, updatedInventory)
	assert.Equal(t, 8, updatedInventory.Quantity)
}

func TestInventoryRepo_GetPurchasedItems(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewInventoryRepo(db)
	_, _ = repo.CreateInventory(1, "item1", 5)
	_, _ = repo.CreateInventory(1, "item2", 3)
	items, err := repo.GetPurchasedItems(1)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestProductRepo_GetItemPrice(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewProductRepo(db)
	_ = db.Create(&models.Product{Name: "product1", Price: 100})
	price, err := repo.GetItemPrice("product1")
	assert.NoError(t, err)
	assert.Equal(t, 100, price)
}

func TestProductRepo_GetItemPrice_NotFound(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewProductRepo(db)

	price, err := repo.GetItemPrice("nonexistent_product")
	assert.Error(t, err)
	assert.Equal(t, 0, price)
}

func TestTransactionRepo_CreateTransaction(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewTransactionRepo(db)
	transaction, err := repo.CreateTransaction(1, 2, 50)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, uint(1), transaction.FromUserID)
	assert.Equal(t, uint(2), transaction.ToUserID)
	assert.Equal(t, 50, transaction.Amount)
}

func TestTransactionRepo_GetGiftsGivenByUser(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewTransactionRepo(db)
	_ = db.Create(&models.Transaction{FromUserID: 1, ToUserID: 2, Amount: 50})
	_ = db.Create(&models.Transaction{FromUserID: 1, ToUserID: 3, Amount: 30})
	gifts, err := repo.GetGiftsGivenByUser(1)
	assert.NoError(t, err)
	assert.Len(t, gifts, 2)
	assert.Equal(t, string("2"), gifts[0].ToUser)
	assert.Equal(t, 50, gifts[0].Amount)
}

func TestTransactionRepo_GetGiftsGivenToUser(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewTransactionRepo(db)
	_ = db.Create(&models.Transaction{FromUserID: 2, ToUserID: 1, Amount: 50})
	_ = db.Create(&models.Transaction{FromUserID: 3, ToUserID: 1, Amount: 30})
	gifts, err := repo.GetGiftsGivenToUser(1)
	assert.NoError(t, err)
	assert.Len(t, gifts, 2)
	assert.Equal(t, string("2"), gifts[0].FromUser)
	assert.Equal(t, 50, gifts[0].Amount)
}

func TestUserRepo_CreateUser(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewUserRepo(db)
	user, err := repo.CreateUser("testuser", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password123", user.Password)
}

func TestUserRepo_GetUserByUsername(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewUserRepo(db)
	_, err := repo.CreateUser("testuser", "password123")
	assert.NoError(t, err)

	user, err := repo.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password123", user.Password)
}

func TestUserRepo_UpdateUser(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewUserRepo(db)
	user, _ := repo.CreateUser("testuser", "password123")

	user.Username = "newusername"
	updatedUser, err := repo.UpdateUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, "newusername", updatedUser.Username)
}

func TestUserRepo_UpdateTwoUsers(t *testing.T) {
	db := getTestDB(t)
	defer clearTables(db)
	applyMigrations(db)

	repo := NewUserRepo(db)
	user1, _ := repo.CreateUser("user1", "password1")
	user2, _ := repo.CreateUser("user2", "password2")

	user1.Username = "updateduser1"
	user2.Username = "updateduser2"
	updatedUser1, updatedUser2, err := repo.UpdateTwoUsers(user1, user2)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser1)
	assert.NotNil(t, updatedUser2)
	assert.Equal(t, "updateduser1", updatedUser1.Username)
	assert.Equal(t, "updateduser2", updatedUser2.Username)
}
