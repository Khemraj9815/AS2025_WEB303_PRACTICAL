package database

import (
	"log"

	"github.com/douglasswmcst/ss2025_web303/practicals/practical5/order-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate order tables
	err = DB.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if err != nil {
		return err
	}

	log.Println("Order-service database connected and migrated")
	return nil
}
