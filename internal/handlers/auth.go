package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/database"
	"github.com/rideaware/admin-panel/internal/middleware"

	"github.com/gin-gonic/gin"
)

// LoginGet renders the login page using the "login.html" template with HTTP 200 status.
func LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// LoginPost handles POST /login form submissions, authenticates the user, creates a session, and redirects to "/" on success.
// On invalid credentials it renders the login page with HTTP 401 and an error message; if session retrieval or saving fails it aborts with HTTP 500.
func LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	admin, err := database.GetAdmin(username)
	if err != nil || !database.VerifyPassword(admin.Password, password) {
		c.HTML(http.StatusUnauthorized, "login.html",
			gin.H{"error": "Invalid username or password"})
		return
	}

	session, err := middleware.GetStore().Get(c.Request, "session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Values["username"] = username
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// Logout invalidates the current user session if one exists and redirects the client to the login page.
// If the session cannot be retrieved, the handler still redirects to "/login".
func Logout(c *gin.Context) {
	session, err := middleware.GetStore().Get(c.Request, "session")
	if err == nil {
		session.Options.MaxAge = -1
		session.Save(c.Request, c.Writer)
	}
	c.Redirect(http.StatusFound, "/login")
}