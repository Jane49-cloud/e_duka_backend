package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
	var err error
	// user := os.Getenv("DB_USER")
	// dbname := os.Getenv("DB_NAME")
	// password := os.Getenv("DB_PASSWORD")
	// host := "localhost"

	db_url := os.Getenv("DB_URL")

	// db_url := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s port=%s host=/var/run/postgresql", user, dbname, sslmode, password, port)
	// db_url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)
	fmt.Printf("db_url \n%s\n", db_url)
	Database, err = gorm.Open(postgres.Open(db_url), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected to the database")
	}
}
