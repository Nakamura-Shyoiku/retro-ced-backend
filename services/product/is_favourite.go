package product

import (
	"errors"

	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

type IsFavourite struct {
	ProductID int64
	GUID      string
	user      models.User
	favourite models.Favourites
}

func NewIsFavourite(ProductID int64, GUID string) *IsFavourite {
	n := new(IsFavourite)
	n.ProductID = ProductID
	n.GUID = GUID
	return n
}

func (f *IsFavourite) Data() bool {
	return f.favourite.Id > 0
}

func (f *IsFavourite) Do() (err error) {
	if err = f.validate(); err != nil {
		return err
	}

	if err = f.getUser(); err != nil {
		return err
	}

	if err = f.getFav(); err != nil {
		return err
	}

	return nil
}

func (f *IsFavourite) validate() (err error) {
	if f.ProductID <= 0 {
		err = errors.New("valid 'product_id' is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (f *IsFavourite) getUser() (err error) {
	f.user, err = models.GetUserByGUID(f.GUID)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		return err
	}
	return nil
}

func (f *IsFavourite) getFav() (err error) {
	f.favourite, err = new(models.Favourites).GetFavourite(f.ProductID, f.user.Id)
	if err != nil {
		log.WithError(err).Error("failed to get favourite")
		return err
	}

	return nil
}
