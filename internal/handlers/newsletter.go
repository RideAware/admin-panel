package handlers

import (
	"net/http"

	"github.com/rideaware/admin-panel/internal/email"

	"github.com/gin-gonic/gin"
)

// SendUpdateGet renders the update form page using the "send_update.html" template and responds with HTTP 200 OK.
func SendUpdateGet(c *gin.Context) {
	c.HTML(http.StatusOK, "send_update.html", gin.H{})
}

// SendUpdatePost handles POST requests to submit a newsletter update.
// It reads "subject" and "body" from the form, calls email.SendUpdate(subject, body),
// and renders the "send_update.html" template with gin.H{"error": message} when sending fails
// or gin.H{"success": message} when sending succeeds, returning HTTP 200 in both cases.
func SendUpdatePost(c *gin.Context) {
	subject := c.PostForm("subject")
	body := c.PostForm("body")

	// validate inputs
	if strings.TrimSpace(subject) == "" || strings.TrimSpace(body) == {
		c.HTML(http,StatusBadRequest, "send_update.html",
			gin.H{"error": "Subject and message cannot be empty"})
		return
	}

	message, err := email.SendUpdate(subject, body)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "send_update.html",
			gin.H{"error": message})
		return
	}

	c.HTML(http.StatusOK, "send_update.html",
		gin.H{"success": message})
}