package database

import (
	"blog-website/models"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=021103 dbname=stocksdb port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	} else {
		fmt.Println("Successfully connected!")
	}

	DB = database

	database.AutoMigrate(
		&models.User{},
		&models.Blog{},
	)
}
