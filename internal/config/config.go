package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	PGHost         string
	PGPort         string
	PGUser         string
	PGPassword     string
	PGDatabase     string
	SMTPServer     string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SenderEmail    string
	AdminUsername  string
	AdminPassword  string
	SecretKey      string
	BaseURL        string
}

var Current *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Port:          getEnv("PORT", "5001"),
		PGHost:        getEnv("PG_HOST", "localhost"),
		PGPort:        getEnv("PG_PORT", "5432"),
		PGUser:        getEnv("PG_USER", "postgres"),
		PGPassword:    getEnv("PG_PASSWORD", ""),
		PGDatabase:    getEnv("PG_DATABASE", "newsletter"),
		SMTPServer:    getEnv("SMTP_SERVER", ""),
		SMTPPort:      getEnvInt("SMTP_PORT", 465),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		SenderEmail:   getEnv("SENDER_EMAIL", ""),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "changeme"),
		SecretKey:     getEnv("SECRET_KEY", "your-secret-key"),
		BaseURL:       getEnv("BASE_URL", "localhost:5001"),
	}

	if cfg.SenderEmail == "" {
		cfg.SenderEmail = cfg.SMTPUser
	}

	Current = cfg
	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
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
		log.Printf("Invalid integer for %s: %v", key, err)
		return defaultValue
	}
	return intVal
}