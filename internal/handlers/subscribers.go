package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/database"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	emails, err := database.GetAllEmails()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "admin_index.html",
		gin.H{"emails": emails})
}