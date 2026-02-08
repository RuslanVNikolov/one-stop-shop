package database

import (
	"database/sql"
	_ "fmt"
	"log"
	"time"

	"github.com/RuslanVNikolov/one-stop-shop/backend/auth-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) *gorm.DB {
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✓ Database connection established")
	return db
}

func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Println("⚠ Warning: Could not create uuid-ossp extension:", err)
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.OAuthAccount{},
	)

	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	log.Println("✓ Database migrations completed successfully")
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetDB(db *gorm.DB) (*sql.DB, error) {
	return db.DB()
}
