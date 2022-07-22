package pages

import (
	"github.com/apex/log"
	"github.com/ulventech/retro-ced-backend/models"

	"github.com/ulventech/retro-ced-backend/scraper/pages/amorevintage"
	"github.com/ulventech/retro-ced-backend/scraper/pages/malleries"
	"github.com/ulventech/retro-ced-backend/scraper/pages/tradesy"
)

type crawler interface {
	Crawl() error
}

// RunScraper starts srapers for all sites listed
func RunScraper() {

	log.Info("starting scraper")

	var enabledSites []models.Site
	err := models.GetDBv2().Where("active = 1").Find(&enabledSites).Error
	if err != nil {
		log.WithError(err).Error("could not retrieve enabled sites")
		return
	}

	crawlers := make(map[string]crawler)

	for _, target := range enabledSites {

		log.WithField("target", target.Name).Info("preparing crawler for site")

		var c crawler

		switch target.Name {

		case models.SiteNameTradesy:
			c, err = tradesy.New(target)
			if err != nil {
				log.WithError(err).Errorf("could not create tradesy crawler")
				continue
			} else {
				crawlers[target.Name] = c
			}

		case models.SiteNameMalleries:

			c, err = malleries.New(target)
			if err != nil {
				log.WithError(err).Errorf("could not create malleries crawler")
				continue
			} else {
				crawlers[target.Name] = c
			}

		case models.SiteNameAmoreVintage:

			c, err = amorevintage.New(target)
			if err != nil {
				log.WithError(err).Errorf("could not create amorevintage crawler")
				continue
			} else {
				crawlers[target.Name] = c
			}

		default:
			log.WithField("name", target.Name).Error("unknown crawl target")
			continue
		}
	}

	for name, target := range crawlers {

		go func(name string, target crawler) {

			log.WithField("name", name).Info("starting crawler")

			err := target.Crawl()
			if err != nil {
				log.WithError(err).Errorf("crawling %v failed", name)
			} else {
				log.Infof("%v crawl succeeded", name)
			}

		}(name, target)
	}

}
