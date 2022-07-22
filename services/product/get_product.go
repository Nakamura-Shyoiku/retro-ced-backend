package product

import (
	"errors"
	
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/models"
)

// Product service
type Product struct {
	id      string
	product *models.Product
}

// NewGetProduct instance
func NewGetProduct(id string) *Product {
	return &Product{
		id:      id,
		product: new(models.Product),
	}
}

// Data returns the final output from Do
func (p *Product) Data() *models.Product {
	return p.product
}

// Do tasks
func (p *Product) Do() (err error) {
	if err = p.validate(); err != nil {
		return err
	}

	if err = p.getProduct(); err != nil {
		return err
	}

	return nil
}

func (p *Product) validate() (err error) {
	if p.id == "" {
		err = errors.New("id is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (p *Product) getProduct() (err error) {
	p.product, err = p.product.GetProduct(p.id)
	if err != nil {
		return err
	}

	return nil
}
