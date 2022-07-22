package malleries

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"github.com/ulventech/retro-ced-backend/models"
)

func (c *Crawler) processItem(address string, category string) error {

	res, err := c.cli.Get(address, nil)
	if err != nil {
		return errors.Wrapf(err, "could not retrieve page (url: %v)", address)
	}

	var raw rawAdPage
	err = c.cli.UnpackHTML(res, &raw)
	if err != nil {
		return errors.Wrapf(err, "could not unpack data (url: %v)", address)
	}

	page := raw.adPage()

	product := models.ProductRecord{
		Guid:         uuid.New().String(),
		SiteId:       int64(c.siteID),
		Title:        page.Title,
		Description:  page.Description,
		Brand:        page.Brand,
		ItemNumber:   page.ItemModel,
		Price:        int64(page.Price),
		Measurements: page.Measurements,
		Color:        page.Color,
		Size:         page.Size,
		ProductURL:   &address,

		// NOTE: dummy URL Guid
		// TODO: check Urls table - this specific URL Guid has a siteID belonging to therealreal
		// (the first with the new approach to handling URLs). Does that matter?
		UrlId: 0,
	}

	if len(page.Images) > 0 {
		product.Img = page.Images[0]
	}

	// categories are 1-to-1, expect 'jewelry' that gets translated to 'accessories'
	if category == "jewelry" {
		product.Category = "accessories"
	} else {
		product.Category = category
	}

	// TODO: retail price

	err = models.GetDBv2().Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "product_url"}}, UpdateAll: true}).Create(&product).Error
	if err != nil {
		return errors.Wrapf(err, "could not process product (url: %v)", address)
	}

	return nil

}
