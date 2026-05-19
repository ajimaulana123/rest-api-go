package database

import (
	"os"
	"strings"

	"be/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	dsn := databaseURL
	if !strings.Contains(dsn, "sslmode=") {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		dsn += sep + "sslmode=require"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Supabase pooler (port 6543)
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Tabel sudah ada di Supabase; AutoMigrate bisa bentrok nama constraint.
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := db.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}
