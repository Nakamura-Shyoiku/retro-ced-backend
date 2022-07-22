package controllers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/ulventech/retro-ced-backend/models"
)

// Cleanup will invoke the Cleanup() from the models package
func Cleanup(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"ok": true})

	go func() {
		log.Info("Cleanup request received")

		err := models.Cleanup()
		if err != nil {
			log.WithError(err).Warn("cleanup failed")
			return
		}

		log.Info("cleanup done")
	}()
}
