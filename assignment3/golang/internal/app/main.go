package app

import (
	"context"
	"log"
	"net/http"
	"time"
	"os"
	"golang/internal/handler"
	"golang/internal/middleware"
	"golang/internal/repository"
	"golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	_postgre := _postgres.NewPGXDialect(ctx, dbConfig)

	repos := repository.NewRepositories(_postgre)
	uc := usecase.NewUserUsecase(repos.UserRepository)
	h := handler.NewUserHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/users/", h.UserByID)

	var root http.Handler = mux
	root = middleware.Logger(root)
	root = middleware.APIKey(root)

	log.Println("Starting the Server.")
	log.Fatal(http.ListenAndServe(":8080", root))
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        getEnv("POSTGRES_HOST", "db"),
		Port:        getEnv("POSTGRES_PORT", "5432"),
		Username:    getEnv("POSTGRES_USER", "postgres"),
		Password:    getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:      getEnv("POSTGRES_DB", "mydb"),
		SSLMode:     getEnv("POSTGRES_SSLMODE", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}