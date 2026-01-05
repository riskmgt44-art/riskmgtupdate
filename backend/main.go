// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"riskmgt/config"
	"riskmgt/database"
	"riskmgt/handlers"
	"riskmgt/routes"
	"riskmgt/utils"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to MongoDB
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize handler collections AFTER successful DB connection
	handlers.InitCollections()

	// Setup router
	router := mux.NewRouter()

	// Global middleware
	router.Use(utils.LoggingMiddleware)
	router.Use(utils.RecoveryMiddleware)
	router.Use(utils.CORSMiddleware)

	// Register all routes
	routes.RegisterRoutes(router)

	// Start server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		log.Printf("RiskMGT Backend starting on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	database.Disconnect()
	log.Println("Server stopped gracefully")
}