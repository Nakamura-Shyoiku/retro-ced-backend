package product

import (
	"strconv"
	
	"github.com/ulventech/retro-ced-backend/models"
)

type SearchProducts struct {
	guid     string
	user     models.User
	products []models.Product
	totalCnt uint64
}

func NewSearchProducts(guid string) *SearchProducts {
	ps := new(SearchProducts)
	ps.guid = guid
	ps.products = make([]models.Product, 0)
	return ps
}

func (ps *SearchProducts) Data() ([]models.Product, uint64) {
	return ps.products, ps.totalCnt
}

func (ps *SearchProducts) Do(search string, offset string) (err error) {
	if err = ps.getUser(); err != nil {
		return err
	}

	if err = ps.searchProducts(search, offset); err != nil {
		return err
	}

	if err = ps.getProductsCountBySearch(search); err != nil {
		return err
	}

	return nil
}

func (ps *SearchProducts) getUser() (err error) {
	ps.user, err = models.GetUserByGUID(ps.guid)
	if err != nil {
		return err
	}

	return nil
}

func (ps *SearchProducts) searchProducts(search string, offset string) (err error) {
	i, err := strconv.Atoi(offset)
	i = i * 18
	s := strconv.Itoa(i)

	ps.products, err = models.SearchProducts(search, s, ps.user.Id)
	if err != nil {
		return err
	}

	return nil
}
func (ps *SearchProducts) getProductsCountBySearch(search string) (err error) {
	ps.totalCnt, err = models.GetProductsCountBySearch(search)
	if err != nil {
		return err
	}
	return nil
}
