package middleware

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func Init() {
	if config.Current.SecretKey == "" {
		panic("SECRET_KEY not set")
	}
	store = sessions.NewCookieStore([]byte(config.Current.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: 0,
	}
}

func GetStore() *sessions.CookieStore {
	return store
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session")
		if err != nil || session.Values["username"] == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}