package database

import (
	"os"

	"github.com/cynicade/todo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:"+os.Getenv("DB_PASS")+"@localhost:5432/postgres"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Todo{})
	DB = db
}
