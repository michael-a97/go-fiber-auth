package database

import (
	"fib/config"
	"fib/pkg/entity"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Config("dbUser"),
		config.Config("dbPassword"),
		config.Config("host"),
		config.Config("dbPort"),
		config.Config("dbName"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failed.\n", err.Error())
	}
	log.Println("Connected to the database.")
	db.AutoMigrate(&entity.User{})

	return db
}
