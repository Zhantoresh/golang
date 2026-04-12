package v1

import (
	"net/http"
	"assignment7/internal/entity"
	"assignment7/internal/usecase"
	"assignment7/utils"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserInterface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
	r := &userRoutes{t}
	
	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)
		
		// Protected routes
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		
		
		protected.GET("/me", r.GetMe)
		
		
		adminOnly := protected.Group("/")
		adminOnly.Use(utils.RoleMiddleware("admin"))
		adminOnly.PATCH("/promote/:id", r.PromoteUser)
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	user := entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     "user", 
	}

	createdUser, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": createdUser})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	user, err := r.t.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"username": user.Username,
		"role": user.Role,
	})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	targetID := c.Param("id")
	
	err := r.t.PromoteUser(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}