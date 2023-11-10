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

	db_url := os.Getenv("DB_URL")
	// db_url := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"

	// dsn := "host=localhost user=postgres password= dbname=eduka_locale port=5432 sslmode=disable"

	Database, err = gorm.Open(postgres.Open(db_url), &gorm.Config{})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to the database")
	}
}
