package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY")) // Ensure this is set in environment variables
var db *pgxpool.Pool
var jwtPool = sync.Pool{
	New: func() interface{} {
		return jwt.New(jwt.SigningMethodHS256)
	},
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func initDB() {
	dbURL := os.Getenv("DATABASE_URL") // Ensure this is set in environment variables

	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Create a new connection pool
	var err error
	db, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Ensure the table exists
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	log.Println("Connected to PostgreSQL")
}

func main() {
	initDB()         // Initialize database
	defer db.Close() // Ensure database connection is closed on exit

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "https://retailer-chi.vercel.app/"},

		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/refresh", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Generate a new token
		expirationTime := time.Now().Add(15 * time.Minute)
		claims.ExpiresAt = expirationTime.Unix()
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = newToken.SignedString(jwtKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user already exists
		var exists bool
		err := db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}

		_, err = db.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var storedPassword string
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond) // Timeout to prevent slow queries
		defer cancel()

		err := db.QueryRow(ctx, "SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT with optimized memory management
		expirationTime := time.Now().Add(15 * time.Minute)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwtPool.Get().(*jwt.Token)
		token.Claims = claims
		tokenString, err := token.SignedString(jwtKey)
		jwtPool.Put(token) // Reuse the token object

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	r.GET("/protected", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Welcome " + claims.Username})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})
	r.Run(":8080")
}
