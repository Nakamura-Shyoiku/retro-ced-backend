package amorevintage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils/client"
	"github.com/ulventech/retro-ced-backend/utils/currency"
)

// New creates a new amorevintage crawler
func New(siteConfig models.Site) (*Crawler, error) {

	client, err := client.NewDefaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create HTTP client")
	}

	rate, err := currency.ConversionRate(currency.LabelJPY, currency.LabelUSD)
	if err != nil {
		// Try again to get the Conversion Rate
		rate, err = currency.ConversionRate(currency.LabelJPY, currency.LabelUSD)
		if err != nil {
			return nil, errors.Wrap(err, "could not determine JPY=>USD conversion rate")
		}
	}

	log.WithField("rate", rate).Info("conversion rate JPY/USD retrieved")

	crawler := Crawler{
		cli:     client,
		maxPage: int(siteConfig.MaxPage),
		siteID:  int(siteConfig.ID),
		sleep:   500 * time.Millisecond,

		jpyToUSDConversionRate: rate,
	}

	return &crawler, nil
}

// Crawl will crawl the actual pages
func (c *Crawler) Crawl() error {

	log.Info("starting amorevintage crawl")

	for _, category := range getCategories() {

		err := c.crawlCategory(category)
		if err != nil {
			return errors.Wrapf(err, "could not crawl amorevintage category (category: %v)", category)
		}

		log.WithField("category", category).Info("successfully crawled amorevingate category")
	}

	log.Info("completed amorevintage category")

	return nil
}

func getCategories() []string {
	return []string{
		"bags",
		"clothing",
		"shoes",
		"accessories",
	}
}

func (c *Crawler) crawlCategory(category string) error {

	address := fmt.Sprintf("%v/%v", amorevintageBaseURL, category)

	for i := 1; i < c.maxPage; i++ {

		params := url.Values{}
		params.Set("page", fmt.Sprint(i))
		params.Set("sort_by", "created-descending")

		res, err := c.cli.Get(address, params)
		if err != nil {
			return errors.Wrapf(err, "could not get search page (category: %v, page: %v)", category, i)
		}

		time.Sleep(1 * time.Second)

		var rawPage searchPage
		err = c.cli.UnpackHTML(res, &rawPage)
		if err != nil {
			return errors.Wrapf(err, "could not unpack search page (category: %v, page: %v)", category, i)
		}

		if len(rawPage.Products) == 0 {
			log.WithField("category", category).WithField("page", i).Info("no more pages to process")
			return nil
		}

		products, err := rawPage.transform()
		if err != nil {
			return errors.Wrapf(err, "could not process search page (category: %v, page: %v)", category, i)
		}

		log.WithField("category", category).WithField("batch", len(products)).WithField("page", i).Info("retrieved amorevintage products")

		for _, product := range products {

			dbProduct := models.ProductRecord{
				Guid:             uuid.New().String(),
				SiteId:           int64(c.siteID),
				Category:         category,
				Title:            product.Name,
				Description:      product.Description,
				Price:            int64(c.getPriceInUSD(product)),
				ItemNumber:       product.SKU,
				ProductCondition: product.getCondition(),
				Accessories:      product.getAccessories(),
				Measurements:     product.getSize(),
				Color:            product.getColor(),
				Size:             product.getSize(),
				Img:              product.ImageURL,
				SubCategory:      product.Category,
				ProductURL:       &product.URL,
			}
			err = models.GetDBv2().Create(&dbProduct).Error
			//err = models.GetDBv2().Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "product_url"}}, UpdateAll: true}).Create(&dbProduct).Error
			if err != nil {
				return errors.Wrapf(err, "could not process product (url: %v)", product.URL)
			}

		}
	}

	return nil
}
