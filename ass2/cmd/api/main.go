package main

import (
	"ass2/internal/handlers"
	"ass2/internal/middleware"
	"log"
	"net/http"

	_ "ass2/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task Management API
// @version 1.0
// @description A simple task management API with authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@taskapi.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY

func main() {
	taskHandler := handlers.NewTaskHandler()

	mux := http.NewServeMux()
	mux.Handle("/v1/tasks", taskHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	handler := middleware.LoggingMiddleware(middleware.AuthMiddleware(mux))

	log.Println("Server starting on :8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
