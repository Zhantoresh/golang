package app

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustConnectDB() *sqlx.DB {
	dsn := "host=localhost port=5432 user=postgres password=A!$ha1973zha dbname=assignment5 sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection failed: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed: ", err)
	}

	log.Println("Database connected")
	return db
}