package controllers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	services "github.com/ulventech/retro-ced-backend/services/analytics"
)

type Analytics struct {
	Base
}

func (b *Analytics) TrackLinkClick(c *gin.Context) {
	link := c.Query("to")
	if link == "" {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "link is invalid"))
		return
	}

	// Decode link
	decodedLink, err := base64.URLEncoding.DecodeString(link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, b.errorMsg(c, err.Error()))
		return
	}

	s := services.NewTrackClick(string(decodedLink))
	if err := s.Do(); err != nil {
		c.JSON(http.StatusInternalServerError, b.errorMsg(c, err.Error()))
		return
	}

	c.Redirect(http.StatusPermanentRedirect, string(decodedLink))
}
