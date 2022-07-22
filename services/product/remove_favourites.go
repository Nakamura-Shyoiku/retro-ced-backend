package product

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type RemoveFavourite struct {
	guid					string
	user					models.User
	favourite 		models.Favourites
}

func NewRemoveFavourite(guid string) *RemoveFavourite {
	f := new(RemoveFavourite)
	f.guid = guid
	return f
}

func (f *RemoveFavourite) Do(productId int64) (err error) {
	if err = f.getUser(); err != nil {
		return err
	}
	if err = f.removeFavourite(productId); err != nil {
		return err
	}

	return nil
}

func (f *RemoveFavourite) removeFavourite(productId int64) error {
	af := &models.Favourites{
		ProductId: 	productId,
		UserId:     f.user.Id,
	}
	return af.RemoveFavourite()
}

func (f *RemoveFavourite) getUser() (err error) {
	f.user, err = models.GetUserByGUID(f.guid)
	if err != nil {
		return err
	}

	return nil
}
