package models

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	FromUserID uint `gorm:"not null"`
	ToUserID   uint `gorm:"not null"`
	Amount     int  `gorm:"not null"`
}
