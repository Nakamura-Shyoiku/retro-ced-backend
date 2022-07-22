package product

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type Brands struct {
	brands []models.Brand
}

func NewBrands() *Brands {
	bs := new(Brands)
	return bs
}

func (bs *Brands) Data() []models.Brand {
	return bs.brands
}

func (bs *Brands) Do() (err error) {
	if err = bs.getBrands(); err != nil {
		return err
	}

	return nil
}

func (bs *Brands) getBrands() (err error) {
	bs.brands, err = models.GetBrands()
	if err != nil {
		return err
	}
	return nil
}