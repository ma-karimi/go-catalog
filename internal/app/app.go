package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-catalog/internal/http"
)

// Run starts the HTTP server
func Run() {
	router := gin.Default()

	// Setup routes
	http.SetupRoutes(router)

	// Start server on port 8080
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}



