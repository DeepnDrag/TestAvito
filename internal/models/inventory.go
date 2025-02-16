package models

type Inventory struct {
	UserID   uint   `gorm:"primaryKey;not null"`
	ItemType string `gorm:"primaryKey;not null"`
	Quantity int    `gorm:"not null"`
}
