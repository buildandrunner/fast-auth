package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mar-cial/space-auth/internal/adapter/handler"
	redisRepo "github.com/mar-cial/space-auth/internal/adapter/repository/redis"
	"github.com/mar-cial/space-auth/internal/core/service"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	redisClient := redis.NewClient(options)

	authRepo := redisRepo.NewRedisAuthRepository(redisClient)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/logout", authHandler.Logout)

	log.Printf("Starting server on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
