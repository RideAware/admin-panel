package database

import (
	"context"
	"time"
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

// Init initializes the package database connection using values from cfg, sets connection pool limits,
// creates required tables, and ensures a default admin user exists.
// It assigns the opened *sql.DB to the package-level db and will terminate the program if establishing
// or verifying the connection fails.
//
// cfg provides PostgreSQL connection parameters and the default admin credentials used to create the
// default admin user when missing.
func Init(cfg *config.Config) {
	password := url.PathEscape(cfg.PGPassword)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection successful!")
	createTables()
	createDefaultAdmin(cfg)
}

// Close closes the package-level database connection if it has been initialized.
// It is safe to call multiple times; if no connection exists, the call is a no-op.
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

// createTables creates the required database tables if they do not already exist.
// It ensures the subscribers, admin_users, and newsletters tables are present; errors
// encountered while creating individual tables are logged but do not abort the process.
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

// GetAllEmails retrieves all subscriber email addresses from the database.
// It returns a slice of email strings and any error encountered while querying or scanning rows.
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

// GetAdmin retrieves the admin user with the given username.
// It returns a pointer to the Admin when a matching row exists. If no admin is found, it returns an error "admin not found"; other database errors are returned unchanged.
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

// createDefaultAdmin ensures a default admin user exists by inserting cfg.AdminUsername
// with a bcrypt-hashed cfg.AdminPassword into the admin_users table; the insert is
// idempotent (no-op if the username already exists). If password hashing fails the
// function terminates the process; insertion errors are logged.
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

// LogNewsletter inserts a newsletter record with the provided subject and body into the newsletters table.
// It returns any error encountered while inserting the record.
func LogNewsletter(subject, body string) error {
	_, err := db.Exec(
		"INSERT INTO newsletters (subject, body) VALUES ($1, $2)",
		subject, body,
	)
	return err
}

// hashPassword generates a bcrypt hash for the given plaintext password.
// It uses bcrypt.DefaultCost and returns the hashed password as a string and any error encountered.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hash), err
}

// VerifyPassword reports whether the provided plaintext password matches the given bcrypt hash.
// It returns true if the password matches, false otherwise.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}