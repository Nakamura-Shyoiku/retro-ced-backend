package product

import (
	"strconv"
	
	"github.com/ulventech/retro-ced-backend/models"
)

type ProductsByBrand struct {
	guid     string
	user     models.User
	products []models.Product
	totalCnt uint64
}

func NewGetProductsByBrand(guid string) *ProductsByBrand {
	ps := new(ProductsByBrand)
	ps.guid = guid
	ps.products = make([]models.Product, 0)
	return ps
}

func (ps *ProductsByBrand) Data() ([]models.Product, uint64) {
	return ps.products, ps.totalCnt
}

func (ps *ProductsByBrand) Do(brand string, offset string) (err error) {
	if err = ps.getUser(); err != nil {
		return err
	}

	if err = ps.getProductsByBrand(brand, offset); err != nil {
		return err
	}

	if err = ps.getProductsCountByBrand(brand); err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByBrand) getUser() (err error) {
	ps.user, err = models.GetUserByGUID(ps.guid)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByBrand) getProductsByBrand(brand string, offset string) (err error) {
	i, err := strconv.Atoi(offset)
	i = i * 18
	s := strconv.Itoa(i)

	ps.products, err = models.GetProductsByBrand(brand, s, ps.user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByBrand) getProductsCountByBrand(brand string) (err error) {
	ps.totalCnt, err = models.GetProductsCountByBrand(brand)
	if err != nil {
		return err
	}
	return nil
}
