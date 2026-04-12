package app

import (
	"log"

	"assignment7/internal/controller/http/v1"
	"assignment7/internal/entity"
	"assignment7/internal/usecase"
	"assignment7/internal/usecase/repo"
	"assignment7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run() {
	
	pg, err := postgres.NewPostgres()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	
	err = pg.Conn.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Migration Error: %v", err)
	}

	
	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)


	handler := gin.Default()
	v1.NewRouter(handler, userUseCase)

	
	log.Println("The server is running successfully on port 8090...")
	if err := handler.Run(":8090"); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}