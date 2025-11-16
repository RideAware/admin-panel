package main

import (
	"log"
	"os"

	"github.com/rideaware/admin-panel/internal/config"
	"github.com/rideaware/admin-panel/internal/database"
	"github.com/rideaware/admin-panel/internal/handlers"
	"github.com/rideaware/admin-panel/internal/middleware"

	"github.com/gin-gonic/gin"
)

// main is the program entry point for the admin panel. It loads configuration,
// initializes middleware and the database (closed on exit), configures a Gin
// router with HTML templates and static assets, registers public and
// authenticated routes, and starts the HTTP server on the configured port.
func main() {
	cfg := config.Load()

	// Set Gin mode based on environment (default to release)
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	middleware.Init()
	database.Init(cfg)
	defer database.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Trust only localhost proxy in production
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost", "::1"})

	router.LoadHTMLGlob("web/templates/*.html")
	router.Static("/static", "web/static")

	router.GET("/login", handlers.LoginGet)
	router.POST("/login", handlers.LoginPost)
	router.GET("/logout", handlers.Logout)

	protected := router.Group("/")
	protected.Use(middleware.Auth())
	protected.GET("/", handlers.IndexGet)
	protected.GET("/send_update", handlers.SendUpdateGet)
	protected.POST("/send_update", handlers.SendUpdatePost)

	log.Printf("Server running on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}