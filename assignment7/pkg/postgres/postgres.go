package postgres

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Postgres struct {
	Conn *gorm.DB
}


func NewPostgres() (*Postgres, error) {
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),     
		os.Getenv("DB_USER"),     
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DB_NAME"),     
		os.Getenv("DB_PORT"),     
	)

	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Couldn't connect to the database: %v", err)
		return nil, err
	}

	log.Println("Successful connection to the assignment7_db database!")

	return &Postgres{Conn: db}, nil
}