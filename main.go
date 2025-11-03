package main

import (
	"log"

	"prenup/config"
	"prenup/handlers"
	"prenup/repository"
	"prenup/services"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	if err := config.CreateOTPTable(); err != nil {
		log.Fatalf("Failed to create OTP table: %v", err)
	}

	otpRepo := repository.NewOTPRepository(config.DB)
	otpService := services.NewOTPService(otpRepo)
	otpHandler := handlers.NewOTPHandler(otpService)

	router := gin.Default()

	router.POST("/otp/request", otpHandler.RequestOTP)
	router.POST("/otp/validate", otpHandler.ValidateOTP)

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
