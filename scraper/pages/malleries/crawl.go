package malleries

import (
	"fmt"
	"net/url"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils/client"
)

// New will create a new malleries crawler.
func New(siteConfig models.Site) (*Crawler, error) {

	client, err := client.NewDefaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create HTTP client")
	}

	crawler := Crawler{
		cli:     client,
		maxPage: int(siteConfig.MaxPage),
		sleep:   time.Duration(int64(time.Millisecond) * siteConfig.Sleep), // TODO: do casting the other way around
		siteID:  siteConfig.ID,
	}

	return &crawler, nil
}

// Crawl will crawl the malleries.com website for 'bags', 'shoes', 'accessories' and 'clothing'
func (c *Crawler) Crawl() error {

	log.Info("starting malleries crawl")

	for _, category := range getCategories() {

		err := c.crawlCategory(category)
		if err != nil {
			return errors.Wrapf(err, "error crawling malleries (category: %v)", category)
		}

		log.WithField("category", category).Info("successfully crawled malleries category")
	}

	log.Info("completed malleries crawl")

	return nil
}

func getCategories() []string {

	return []string{
		"bags",
		"shoes",
		"accessories",
		"clothing",
		"jewelry",
	}
}

func (c *Crawler) crawlCategory(category string) error {

	log.WithField("category", category).
		Info("crawling malleries category")

	for i := 1; i <= c.maxPage; i++ {

		// TODO: once stable enough, this can also perhaps be logged with lowered severity
		log.WithField("page", i).WithField("category", category).Info("crawling malleries page")

		params := getParamsForCategorySearch(category)
		params.Set("page", fmt.Sprint(i))

		// TODO: proper throttling
		time.Sleep(c.sleep)

		res, err := c.cli.Get(malleriesEndpoint, params)
		if err != nil {
			return errors.Wrapf(err, "could not get malleries items on page (url: %v, page: %v)", malleriesEndpoint, i)
		}

		var page resultsPage
		err = c.cli.UnpackHTML(res, &page)
		if err != nil {
			return errors.Wrapf(err, "could not unpack malleries items on page (url: %v, page: %v)", malleriesEndpoint, i)
		}

		// TODO: improve pagination
		if len(page.Items) == 0 {
			log.WithField("category", category).
				Info("no more malleries items found")
			break
		}

		for _, item := range page.Items {
			err = c.processItem(item, category)
			if err != nil {
				log.WithField("url", item).WithField("category", category).WithError(err).Error("could not process item")
				// TODO: once stable enough, continue on errors
			}
		}

		log.WithField("page", i).
			WithField("category", category).
			WithField("count", len(page.Items)).
			Info("retrieved malleries items")
	}

	return nil
}

func getParamsForCategorySearch(category string) url.Values {

	params := url.Values{}

	switch category {
	case "clothing":
		params.Set("category", "clothing")
		params.Set("gender", "women")
	case "bags":
		params.Set("category", "handbags")
	case "shoes":
		params.Set("category", "shoes")
		params.Set("gender", "women")
	case "accessories":
		params.Set("type", "accessories")
		params.Set("gender", "female")
	case "jewelry":
		params.Set("category", "jewelry")
		params.Set("gender", "unisex|women")
	}

	return params
}
