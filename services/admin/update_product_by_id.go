package admin

import (
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

type UpdateProductById struct {
	Guid        string  `json:"guid"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Price       int64  `json:"price"`
	RetailPrice int64  `json:"retail_price"`
	UrlId       int64  `json:"url_id"`
	Url         string `json:"url"`
	Color       string `json:"color"`
	Size        string `json:"size"`
	ShoeSize    string `json:"shoe_size"`
	SubCategory string `json:"sub_category"`
	Featured    string `json:"featured"`
	Approved    bool   `json:"approved"`
	product     models.Product
}

func NewUpdateProductById(data UpdateProductById) *UpdateProductById {
	return &UpdateProductById{
		Guid:        data.Guid,
		Title:       data.Title,
		Category:    data.Category,
		Description: data.Description,
		Brand:       data.Brand,
		Price:       data.Price,
		RetailPrice: data.RetailPrice,
		Color:       data.Color,
		Size:        data.Size,
		ShoeSize:    data.ShoeSize,
		SubCategory: data.SubCategory,
		Featured:    data.Featured,
		Approved:    data.Approved,
		UrlId:       data.UrlId,
		Url:         data.Url,
	}
}

func (u *UpdateProductById) Do() (err error) {
	if err = u.validate(); err != nil {
		return err
	}

	// only update if it's not a pseudo URL
	if u.UrlId != 0 {
		if err = u.updateUrl(); err != nil {
			return err
		}
	}

	if err = u.updateProductById(); err != nil {
		return err
	}

	if err = u.getUpdated(); err != nil {
		return err
	}

	return
}

func (u *UpdateProductById) validate() (err error) {
	return nil
}

func (u *UpdateProductById) updateUrl() error {
	err := new(models.Url).UpdateUrl(u.UrlId, u.Url)
	if err != nil {
		log.WithError(err).Error("failed to update url")
		return err
	}

	return nil
}

func (u *UpdateProductById) updateProductById() error {
	// update pseudoURL if needed

	err := models.UpdateProductByGuid(
		u.Guid,
		u.Title,
		u.Description,
		u.Price,
		u.RetailPrice,
		u.Color,
		u.Size,
		u.ShoeSize,
		u.SubCategory,
		u.Featured,
		u.Approved,
	)
	if err != nil {
		log.WithError(err).Error("failed to update product")
		return err
	}

	return nil
}

func (u *UpdateProductById) getUpdated() (err error) {
	u.product, err = models.GetProductByGuid(u.Guid)
	if err != nil {
		log.WithError(err).Error("failed to get updated product")
		return err
	}
	return nil
}
