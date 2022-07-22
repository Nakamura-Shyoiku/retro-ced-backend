package product

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type AddFavourite struct {
	guid					string
	user					models.User
	favourite 		models.Favourites
}

func NewAddFavourite(guid string) *AddFavourite {
	f := new(AddFavourite)
	f.guid = guid
	return f
}

func (f *AddFavourite) Do(productId int64) (err error) {
	if err = f.getUser(); err != nil {
		return err
	}
	if err = f.addFavourite(productId); err != nil {
		return err
	}

	return nil
}

func (f *AddFavourite) addFavourite(productId int64) error {
	af := &models.Favourites{
		ProductId: 	productId,
		UserId:     f.user.Id,
	}
	return af.AddFavourite()
}

func (f *AddFavourite) getUser() (err error) {
	f.user, err = models.GetUserByGUID(f.guid)
	if err != nil {
		return err
	}

	return nil
}
