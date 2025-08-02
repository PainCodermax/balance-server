package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/paincodermax/balance-service/internal/api"
	"github.com/paincodermax/balance-service/internal/database"
	"github.com/paincodermax/balance-service/internal/expense"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI must be set in your .env file")
	}
	db := database.ConnectDB(mongoURI)
	expenseService := expense.NewService(db)
	expenseHandler := api.NewHandler(expenseService)

	router := api.SetupRouter(expenseHandler)
	serverAddr := ":" + port
	log.Println("Starting server on port: ", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
