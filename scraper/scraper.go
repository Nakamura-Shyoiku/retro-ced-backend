package scraper

import (
	"fmt"

	"github.com/anaskhan96/soup"
	"github.com/apex/log"
	raven "github.com/getsentry/raven-go"

	"github.com/ulventech/retro-ced-backend/scraper/pages"
	"github.com/ulventech/retro-ced-backend/utils/env"
)

// Start scraper
func Start() {
	if env.IsDev() {
		log.Info("soup debug mode: true")
		soup.SetDebug(true)
	} else {
		log.Info("soup debug mode: false")
		soup.SetDebug(false)
	}

	// Recover so we don't crash the entire server
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("PANIC!!! Recovered in scraper Start: %s", r)
			raven.CaptureErrorAndWait(err, nil)
			log.Errorf(err.Error())
		}
	}()

	// Run all the scrapers
	pages.RunScraper()
}
