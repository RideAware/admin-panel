package middleware

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

// Init initializes the package-level cookie store used for session management.
// It panics if config.Current.SecretKey is empty.
// The created store is configured with Path "/", MaxAge one week, HttpOnly true, Secure false, and SameSite 0.
func Init() {
	if config.Current == nil {
		panic("config was not loaded; call config.Load() before middleware.Init()")
	}

	if config.Current.SecretKey == "" {
		panic("SECRET_KEY not set")
	}
	store = sessions.NewCookieStore([]byte(config.Current.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

// GetStore returns the package-level Gorilla cookie store used for session management.
// It may be nil if Init has not been called.
func GetStore() *sessions.CookieStore {
	return store
}

// Auth enforces session-based authentication for Gin handlers.
// If the request has no session named "session" or the session lacks a "username" value,
// the middleware redirects to "/login" (HTTP 302) and aborts further handling.
// Otherwise the middleware calls the next handler in the chain.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil {
			c.String(http.StatusInternalServerError, "Session store not initialized.")
			c.Abort()
			return
		}
		session, err := store.Get(c.Request, "session")
		if err != nil || session.Values["username"] == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}