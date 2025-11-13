package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/database"
	"github.com/rideaware/admin-panel/internal/middleware"

	"github.com/gin-gonic/gin"
)

func LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

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

func Logout(c *gin.Context) {
	session, err := middleware.GetStore().Get(c.Request, "session")
	if err == nil {
		session.Options.MaxAge = -1
		session.Save(c.Request, c.Writer)
	}
	c.Redirect(http.StatusFound, "/login")
}