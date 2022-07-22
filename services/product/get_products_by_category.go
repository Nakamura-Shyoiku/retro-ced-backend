package product

import (
	"strconv"
	
	"github.com/ulventech/retro-ced-backend/models"
)

type ProductsByCategory struct {
	guid     string
	user     models.User
	products []models.Product
	totalCnt uint64
}

func NewGetProductsByCategory(guid string) *ProductsByCategory {
	ps := new(ProductsByCategory)
	ps.guid = guid
	ps.products = make([]models.Product, 0)
	return ps
}

func (ps *ProductsByCategory) Data() ([]models.Product, uint64) {
	return ps.products, ps.totalCnt
}

func (ps *ProductsByCategory) Do(category string, offset string, filterCategory []string, filterBrand []string, filterColor []string, filterSize []string, filterShoeSize []string) (err error) {
	if err = ps.getUser(); err != nil {
		return err
	}

	if err = ps.getProductsByCategory(category, offset, filterCategory, filterBrand, filterColor, filterSize, filterShoeSize); err != nil {
		return err
	}

	if err = ps.getProductsCountByCategory(category, filterCategory, filterBrand, filterColor, filterSize, filterShoeSize); err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByCategory) getUser() (err error) {
	ps.user, err = models.GetUserByGUID(ps.guid)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByCategory) getProductsByCategory(category string, offset string, filterCategory []string, filterBrand []string, filterColor []string, filterSize []string, filterShoeSize []string) (err error) {
	i, err := strconv.Atoi(offset)
	i = i * 18
	s := strconv.Itoa(i)

	ps.products, err = models.GetProductsByCategory(category, s, filterCategory, filterBrand, filterColor, filterSize, filterShoeSize, ps.user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ProductsByCategory) getProductsCountByCategory(category string, filterCategory []string, filterBrand []string, filterColor []string, filterSize []string, filterShoeSize []string) (err error) {
	ps.totalCnt, err = models.GetProductsCountByCategory(category, filterCategory, filterBrand, filterColor, filterSize, filterShoeSize)
	if err != nil {
		return err
	}

	return nil
}
