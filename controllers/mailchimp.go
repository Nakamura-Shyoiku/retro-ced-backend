package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ulventech/retro-ced-backend/services/mailchimp"
)

// Mailchimp controller
type Mailchimp struct {
	Base
}

// AddSubscriber adds subscriber to mailchimp list
func (b *Mailchimp) AddSubscriber(c *gin.Context) {
	s := mailchimp.NewAddSubscriber(c.PostForm("email"))
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
