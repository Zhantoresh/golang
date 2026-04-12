package v1

import (
	"time"
	"assignment7/internal/usecase"
	"assignment7/utils"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, u usecase.UserInterface) {
	
	handler.Use(utils.RateLimiterMiddleware(3, time.Minute))

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, u)
	}
}