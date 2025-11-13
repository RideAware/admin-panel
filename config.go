package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string

	// Database
	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDatabase string

	// SMTP
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SenderEmail  string

	// Admin
	AdminUsername string
	AdminPassword string

	// App
	SecretKey string
	BaseURL   string
}

var config *Config

func loadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Debug: Print raw values before any processing
	rawPassword := os.Getenv("PG_PASSWORD")
	log.Printf("Raw PG_PASSWORD length: %d, value: [%s]", len(rawPassword), rawPassword)
	log.Printf("Raw PG_USER: [%s]", os.Getenv("PG_USER"))
	log.Printf("Raw PG_HOST: [%s]", os.Getenv("PG_HOST"))

	cfg := &Config{
		Port:           getEnv("PORT", "5001"),
		PGHost:         getEnv("PG_HOST", "localhost"),
		PGPort:         getEnv("PG_PORT", "5432"),
		PGUser:         getEnv("PG_USER", "postgres"),
		PGPassword:     getEnv("PG_PASSWORD", ""),
		PGDatabase:     getEnv("PG_DATABASE", "newsletter"),
		SMTPServer:     getEnv("SMTP_SERVER", ""),
		SMTPPort:       getEnvInt("SMTP_PORT", 465),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SenderEmail:    getEnv("SENDER_EMAIL", ""),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "changeme"),
		SecretKey:      getEnv("SECRET_KEY", "your-secret-key"),
		BaseURL:        getEnv("BASE_URL", "localhost:5001"),
	}

	// Debug output
	log.Printf("=== Config Loaded ===")
	log.Printf("PG_HOST: %s", cfg.PGHost)
	log.Printf("PG_PORT: %s", cfg.PGPort)
	log.Printf("PG_USER: %s", cfg.PGUser)
	log.Printf("PG_DATABASE: %s", cfg.PGDatabase)
	log.Printf("PG_PASSWORD length: %d", len(cfg.PGPassword))
	log.Printf("BASE_URL: %s", cfg.BaseURL)
	log.Printf("====================")

	if cfg.SenderEmail == "" {
		cfg.SenderEmail = cfg.SMTPUser
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Env var %s not found, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid integer for %s: %v, using default", key, err)
		return defaultValue
	}
	return intVal
}

func maskPassword(pwd string) string {
	if len(pwd) <= 2 {
		return "***"
	}
	return pwd[:2] + "***" + pwd[len(pwd)-2:]
}