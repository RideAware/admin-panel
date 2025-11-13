package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/database"

	"github.com/gin-gonic/gin"
)

// IndexGet handles requests for the admin index page by retrieving all subscriber emails
// and rendering the "admin_index.html" template with those emails.
// If retrieving emails fails, it aborts the request with HTTP 500 and the error.
func IndexGet(c *gin.Context) {
	emails, err := database.GetAllEmails()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_index.html",
		gin.H{"emails": emails})
}