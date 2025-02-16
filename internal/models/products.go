package models

type Product struct {
	Name  string `gorm:"primaryKey;not null"`
	Price int    `gorm:"not null"`
}
