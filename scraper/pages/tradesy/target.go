package tradesy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils/client"
)

const (
	tradesyURL     = "https://www.tradesy.com"
	resultsPerPage = 191

	// for converting relative URLs to absolute
	tradesyScheme   = "https"
	tradesyHostname = "tradesy.com"
)

func getCategories() []string {
	return []string{
		"bags",
		"shoes",
		"clothing",
		"accessories",
	}
}

// Tradesy is a tradesy site crawler
type Tradesy struct {
	client *client.Client

	maxPage int
	sleep   int64
	siteID  uint64
}

// New will create a new tradesy crawler
func New(siteConfig models.Site) (*Tradesy, error) {

	// HTTP client with proxy configuration
	cli, err := client.NewProxyClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create proxy client")
	}

	tradesy := Tradesy{
		client: cli,

		maxPage: int(siteConfig.MaxPage),
		sleep:   siteConfig.Sleep,
		siteID:  siteConfig.ID,
	}

	return &tradesy, nil
}

// Fetch downloads a single ad given its URL
func (t *Tradesy) Fetch(address string) (interface{}, error) {

	res, err := t.client.Get(address, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve page (url: %v)", address)
	}
	var page adPage
	err = t.client.UnpackHTML(res, &page)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unpack data (url: %v)", address)
	}

	var ad adMetadata
	err = json.Unmarshal([]byte(page.MetadataJSON), &ad)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unpack ad metadata for url %s metadata %s", address, page.MetadataJSON)
	}

	payload, err := json.Marshal(ad)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal payload (url: %v)", address)
	}

	return payload, nil
}

// Crawl will crawl the tradesy site for 'bags', 'shoes', 'accessories' and 'clothing' categories
func (t *Tradesy) Crawl() error {

	log.Info("starting tradesy crawl")

	for _, category := range getCategories() {

		// Q: on errors - try to continue or be loud and fail
		// A: let's try to be rigid at first, later we can switch
		err := t.crawlCategory(category)
		if err != nil {
			return errors.Wrapf(err, "error crawling tradesy (category: %v)", category)
		}

		log.WithField("category", category).Info("successfully crawled tradesy category")
	}

	log.Info("completed tradesy crawl")

	return nil
}

func (t *Tradesy) crawlCategory(category string) error {

	log.WithField("category", category).
		Info("crawling tradesy category")

	for i := 1; i < t.maxPage; i++ {

		// TODO: once stable enough, this can also perhaps be logged with lowered severity
		log.WithField("page", i).
			WithField("category", category).
			Info("crawling tradesy page")

		if t.sleep > 0 {
			log.WithField("target", "tradesy").Debugf("sleeping for %d ms", t.sleep)
			time.Sleep(time.Duration(t.sleep) * time.Millisecond)
		}

		address := fmt.Sprintf("%v/%v/", tradesyURL, category)
		params := url.Values{}
		params.Set("page", fmt.Sprint(i))
		params.Set("num_per_page", fmt.Sprint(resultsPerPage))

		res, err := t.client.Get(address, params)
		if err != nil {
			log.WithError(err).Warnf("could not GET HTML on page (url: %v, page: %v)", address, i)
			continue
			// return errors.Wrapf(err, "could not get tradesy items on page (url: %v, page: %v)", address, i)
		}

		var page resultPage
		err = t.client.UnpackHTML(res, &page)
		if err != nil {
			log.WithError(err).Warnf("could not UNPACK HTML on page (url: %v, page: %v)", address, i)
			continue
			//			return errors.Wrapf(err, "could not unpack tradesy search results (url: %v, page: %v)", address, i)
		}

		log.WithField("page", i).
			WithField("category", category).
			WithField("count", len(page.Items)).
			Info("retrieved tradesy items")

		for _, item := range page.Items {

			if item == "" {
				continue
			}

			parsedURL, err := url.Parse(item)
			if err != nil {
				log.WithError(err).
					WithField("url", item).
					Warn("could not process tradesy URL")
				continue
			}

			if !parsedURL.IsAbs() {
				parsedURL.Scheme = tradesyScheme
				parsedURL.Host = tradesyHostname
			}

			// NOTE: continuing despite errors
			err = t.processURL(parsedURL.String(), category)
			if err != nil {
				log.WithError(err).
					WithField("category", category).
					WithField("url", item).
					Warn("could not process tradesy URL")
			}
		}
	}

	return nil
}

func (t *Tradesy) processURL(address string, category string) error {

	log.WithField("url", address).
		WithField("category", category).
		Debug("processing tradesy URL")

	// save URL in the database
	// TODO: this can be part of the products table
	u := &models.Url{}
	err := u.AddUrl(t.siteID, address, category)
	if err != nil {
		return errors.Wrap(err, "could not save tradesy URL in DB")
	}

	if t.sleep > 0 {
		log.WithField("target", "tradesy").Debugf("sleeping for %d ms", t.sleep)
		time.Sleep(time.Duration(t.sleep) * time.Millisecond)
	}

	data, err := t.client.Get(address, nil)
	if err != nil {
		return errors.Wrapf(err, "could not retrieve ad (url: %v)", address)
	}

	var page adPage
	err = t.client.UnpackHTML(data, &page)
	if err != nil {
		return errors.Wrapf(err, "could not unpack ad page (url: %v)", address)
	}

	err = t.out(u, &page, category)
	if err != nil {
		return errors.Wrapf(err, "could not save product (url: %v)", address)
	}

	return nil
}

// TODO: Approved? currently not
func (t *Tradesy) out(linkedURL *models.Url, page *adPage, category string) error {

	var ad adMetadata
	err := json.Unmarshal([]byte(page.MetadataJSON), &ad)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal ad metadata")
	}

	p := &models.Product{

		SiteId:     int64(t.siteID),
		UrlId:      int64(linkedURL.Id),
		Url:        linkedURL.Url,
		ProductURL: linkedURL.Url,
		Category:   category,

		Brand:            strings.ToLower(ad.Brand),
		Title:            ad.Title,
		Description:      ad.Description,
		Price:            int64(ad.Price),
		ProductCondition: ad.Condition,
		Color:            ad.Color,
		Size:             ad.Measurements,
		Measurements:     ad.Measurements,

		Img: page.FeaturedImage,
	}

	if category == "shoes" {
		p.ShoeSize = ad.Size
	}

	if p.Img == "" && len(page.Images) > 0 {
		p.Img = page.Images[0]
	}

	err = p.Create()
	if err != nil {
		return errors.Wrap(err, "could not save product")
	}

	log.WithField("id", p.Guid).Debug("saved tradesy product")

	return nil
}
