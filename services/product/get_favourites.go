package product

import (
	"github.com/ulventech/retro-ced-backend/models"
)

type Favourites struct {
	guid       string
	user       models.User
	favourites []models.Product
	totalCnt   uint64
}

func NewFavourites(guid string) *Favourites {
	fs := new(Favourites)
	fs.guid = guid
	return fs
}

func (fs *Favourites) Data() ([]models.Product, uint64) {
	return fs.favourites, fs.totalCnt
}

func (fs *Favourites) Do() (err error) {
	if err = fs.getUser(); err != nil {
		return err
	}

	if err = fs.getFavourites(); err != nil {
		return err
	}

	if err = fs.getFavouritesCount(); err != nil {
		return err
	}

	return nil
}

func (fs *Favourites) getUser() (err error) {
	fs.user, err = models.GetUserByGUID(fs.guid)
	if err != nil {
		return err
	}

	return nil
}

func (fs *Favourites) getFavourites() (err error) {
	fs.favourites, err = models.GetFavourites(fs.user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (fs *Favourites) getFavouritesCount() (err error) {
	fs.totalCnt, err = models.GetFavouritesCount(fs.user.Id)
	if err != nil {
		return err
	}

	return nil
}
