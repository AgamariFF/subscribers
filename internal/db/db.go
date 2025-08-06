package db

import (
	"subscribers/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(dsn string) *gorm.DB {
	logger.SugaredLogger.Info("Connecting to the database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.SugaredLogger.Errorf("Couldn't connect to PostgreSQL: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.SugaredLogger.Errorf("Couldn't get sql.DB: %v", err)
		return nil
	}
	if err = sqlDB.Ping(); err != nil {
		logger.SugaredLogger.Errorf("Ping to the database failed: %v", err)
		return nil
	}

	logger.SugaredLogger.Info("GORM is connected to PostgreSQL")
	return db
}
