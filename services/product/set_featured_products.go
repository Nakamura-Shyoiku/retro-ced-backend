package product

import (

	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

type SetFeaturedProducts struct {
	GUID     string `json:"guid"`
	Category string `json:"category"`
}

func NewSetFeaturedProducts(data SetFeaturedProducts) *SetFeaturedProducts {
	n := new(SetFeaturedProducts)
	n.GUID = data.GUID
	n.Category = data.Category
	return n
}

func (sf *SetFeaturedProducts) Do() (err error) {
	if err = sf.validate(); err != nil {
		return err
	}

	if err = sf.updateSite(); err != nil {
		return err
	}

	return nil
}

func (sf *SetFeaturedProducts) validate() (err error) {

	return nil
}

func (sf *SetFeaturedProducts) updateSite() error {
	log.Infof("Will update: %v", sf)
	return new(models.Product).SetFeaturedProducts(
		sf.GUID,
		sf.Category,
	)
}
