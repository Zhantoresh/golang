package app

import (
	"context"
	"log"
	"net/http"
	"time"

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

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", root))
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "A!$ha1973zha", 
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}