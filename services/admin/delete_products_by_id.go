package admin

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type DeleteProductsById struct {
	ID []uint64 `json:"Guid"`
}

// FIXME: this is just dumb and again pure overhead

func NewDeleteProductsById(data DeleteProductsById) *DeleteProductsById {
	n := new(DeleteProductsById)
	n.ID = data.ID
	return n
}

func (d *DeleteProductsById) Do() error {
	return models.DeleteProductsById(d.ID)
}
