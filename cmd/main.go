package main

import (
	"TestAvito/internal/config"
	"TestAvito/internal/database"
	"TestAvito/internal/logger"
	"TestAvito/internal/storage"
	"TestAvito/internal/web"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return err
	}

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		return err
	}

	db, err := database.Connection(cfg.Database)
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error(err.Error())
			return
		}
		err = sqlDB.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}()

	err = database.RunMigrations(db)
	if err != nil {
		logger.Error("database run migrations", err.Error())
		return err
	}

	st := storage.New(db)

	server, err := web.New(cfg.Server, cfg.JWT, logger, st)
	if err != nil {
		return err
	}

	return server.Serve()
}
