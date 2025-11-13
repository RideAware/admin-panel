package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/rideaware/admin-panel/internal/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Admin struct {
	Username string
	Password string
}

func Init(cfg *config.Config) {
	password := url.QueryEscape(cfg.PGPassword)

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		cfg.PGUser, password, cfg.PGHost, cfg.PGPort, cfg.PGDatabase,
	)

	log.Printf("Connecting to database: postgres://%s:***@%s:%s/%s",
		cfg.PGUser, cfg.PGHost, cfg.PGPort, cfg.PGDatabase)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection successful!")
	createTables()
	createDefaultAdmin(cfg)
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS subscribers (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS admin_users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS newsletters (
			id SERIAL PRIMARY KEY,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}

	log.Println("Database tables ready.")
}

func GetAllEmails() ([]string, error) {
	rows, err := db.Query("SELECT email FROM subscribers")
	if err != nil {
		log.Printf("Error retrieving emails: %v", err)
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, rows.Err()
}

func GetAdmin(username string) (*Admin, error) {
	var admin Admin
	err := db.QueryRow(
		"SELECT username, password FROM admin_users WHERE username=$1",
		username,
	).Scan(&admin.Username, &admin.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, err
	}

	return &admin, nil
}

func createDefaultAdmin(cfg *config.Config) {
	hashedPassword, err := hashPassword(cfg.AdminPassword)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	_, err = db.Exec(
		"INSERT INTO admin_users (username, password) VALUES ($1, $2) "+
			"ON CONFLICT (username) DO NOTHING",
		cfg.AdminUsername, hashedPassword,
	)
	if err != nil {
		log.Printf("Error creating default admin: %v", err)
	} else {
		log.Println("Default admin user ready.")
	}
}

func LogNewsletter(subject, body string) error {
	_, err := db.Exec(
		"INSERT INTO newsletters (subject, body) VALUES ($1, $2)",
		subject, body,
	)
	return err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hash), err
}

func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}