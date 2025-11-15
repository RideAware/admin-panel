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

// Load loads configuration from environment variables or a .env file and initializes the package-level Current configuration.
// It constructs a Config with sensible defaults for server, PostgreSQL, SMTP, admin credentials, secret key, and base URL.
// If SENDER_EMAIL is not set, it falls back to SMTP_USER. The created Config is assigned to Current and returned.
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
		SecretKey:     getEnv("SECRET_KEY", ""),
		BaseURL:       getEnv("BASE_URL", "localhost:5001"),
	}

	if cfg.SecretKey == "" {
		log.Fatal("SECRET_KEY environment variable must be set!")
	}

	if cfg.SenderEmail == "" {
		cfg.SenderEmail = cfg.SMTPUser
	}

	Current = cfg
	return cfg
}

// getEnv returns the value of the environment variable named by key, or defaultValue if that variable is not set or is empty.
// If the environment variable exists but is the empty string, defaultValue is returned.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt retrieves the environment variable named by key and returns its integer value or defaultValue.
// If the variable is not set, it returns defaultValue. If the variable is set but cannot be parsed as an integer,
// it logs the parse error and returns defaultValue.
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