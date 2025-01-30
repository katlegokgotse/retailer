// main.go
package main

import (
	"net/http"
	"retailer-backend/entities"
	"retailer-backend/interface_adapters/controllers"
	"retailer-backend/interface_adapters/repositories"
	"retailer-backend/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

func main() {
	// Initialize the repository
	userRepo := repositories.NewInMemoryUserRepository()

	// Initialize the use cases
	registerUser := &models.RegisterUser{UserRepo: userRepo}
	loginUser := &models.LoginUser{UserRepo: userRepo, JWTKey: jwtKey}

	// Initialize the controller
	authController := controllers.NewAuthController(registerUser, loginUser)

	// Set up the Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// Protected route
	r.GET("/protected", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			return
		}

		claims := &entities.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Welcome " + claims.Username})
	})

	// Start the server
	r.Run(":8080")
}
