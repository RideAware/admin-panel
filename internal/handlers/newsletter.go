package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/email"

	"github.com/gin-gonic/gin"
)

func SendUpdateGet(c *gin.Context) {
	c.HTML(http.StatusOK, "send_update.html", gin.H{})
}

func SendUpdatePost(c *gin.Context) {
	subject := c.PostForm("subject")
	body := c.PostForm("body")

	message, err := email.SendUpdate(subject, body)
	if err != nil {
		c.HTML(http.StatusOK, "send_update.html",
			gin.H{"error": message})
		return
	}

	c.HTML(http.StatusOK, "send_update.html",
		gin.H{"success": message})
}