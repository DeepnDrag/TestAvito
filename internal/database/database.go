package database

import (
	"TestAvito/internal/config"
	"TestAvito/internal/models"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func Connection(config config.Database) (*gorm.DB, error) {
	log.Println(config)
	if config.User == "" || config.Password == "" || config.Host == "" || config.Port == 0 || config.Name == "" {
		return nil, fmt.Errorf("invalid database configuration")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
	)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	var db *gorm.DB
	var err error
	for attempts := 0; attempts < 3; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("gorm open error after retries: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(models.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(models.Transaction{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(models.Inventory{})
	if err != nil {
		return err
	}

	if !db.Migrator().HasTable(&models.Product{}) {
		err := db.AutoMigrate(&models.Product{})
		if err != nil {
			return err
		}

		products := []models.Product{
			{Name: "t-shirt", Price: 80},
			{Name: "cup", Price: 20},
			{Name: "book", Price: 50},
			{Name: "pen", Price: 10},
			{Name: "powerbank", Price: 200},
			{Name: "hoody", Price: 300},
			{Name: "umbrella", Price: 200},
			{Name: "socks", Price: 10},
			{Name: "wallet", Price: 50},
			{Name: "pink-hoody", Price: 500},
		}

		for _, product := range products {
			err := db.Create(&product).Error
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return fmt.Errorf("db automigrate error: %w", err)
	}
	return nil
}
