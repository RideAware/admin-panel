package main

import (
	"log"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func initSessions() {
	if config.SecretKey == "" {
		log.Fatal("SECRET_KEY not set in configuration")
	}
	store = sessions.NewCookieStore([]byte(config.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: 0,
	}
	log.Println("Sessions initialized")
}