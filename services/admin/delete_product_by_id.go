package admin

import (
	"github.com/ulventech/retro-ced-backend/models"
)

// FIXME: this is just pure overhead

type DeleteProductById struct {
	ID      int64
	product models.Product
}

func NewDeleteProductById(ID int64) *DeleteProductById {
	return &DeleteProductById{
		ID: ID,
	}
}

func (d *DeleteProductById) Do() (err error) {
	return models.DeleteProductById(d.ID)
}
