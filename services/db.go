package services

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connect database

func OpenDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
