package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ulventech/retro-ced-backend/services/trending"
)

// Trending controller
type Trending struct {
	Base
}

// GetPosts returns medium posts
func (b *Trending) GetPosts(c *gin.Context) {
	s := trending.NewGetPosts()
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}
