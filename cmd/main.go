package main

import (
	"expense-tracker/config"
	"expense-tracker/middleware"
	"expense-tracker/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	routes.SetupRoutes(r)

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
