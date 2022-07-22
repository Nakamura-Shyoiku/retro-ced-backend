package controllers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"

	"github.com/ulventech/retro-ced-backend/scraper"
	"github.com/ulventech/retro-ced-backend/scraper/ftp"
)

// Scraper controller
type Scraper struct {
	Base
}

// RunScrapers runs scrapers
func (b *Scraper) RunScrapers(c *gin.Context) {
	log.Info("Dispatching scraper job")
	// TODO: setup a proper queue
	go scraper.Start()

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// RunSFTPScraper runs sftp scraper
func (b *Scraper) RunSFTPScraper(c *gin.Context) {
	log.Info("Dispatching sftp scraper job")
	go ftp.DownloadCJ()

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
