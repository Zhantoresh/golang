package _postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"golang/pkg/modules"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGXDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	log.Println("Waiting for database healthcheck...")

	var db *sqlx.DB
	var err error

	for i := 1; i <= 30; i++ {
		select {
		case <-ctx.Done():
			panic("db connect canceled")
		default:
		}

		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
    		pingErr := db.Ping()
    		if pingErr == nil {
        	log.Println("Database is ready.")
        	AutoMigrate(cfg)
        	return &Dialect{DB: db}
		}
    _ = db.Close()
    err = pingErr
}

		log.Printf("DB not ready (try %d/30): %v\n", i, err)
		time.Sleep(2 * time.Second)
	}

	panic(fmt.Sprintf("database not ready after retries: %v", err))
}

func AutoMigrate(cfg *modules.PostgreConfig) {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}