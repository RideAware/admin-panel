package main

import (
	"log"

	"github.com/rideaware/admin-panel/internal/config"
	"github.com/rideaware/admin-panel/internal/database"
	"github.com/rideaware/admin-panel/internal/handlers"
	"github.com/rideaware/admin-panel/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	middleware.Init()
	database.Init(cfg)
	defer database.Close()

	router := gin.Default()

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