package product

import (
	"errors"
	"time"

	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

// Products service
type Products struct {
	guid        string
	searchQuery string
	category    string
	offset      int
	limit       int
	siteId      string
	featured    string
	sortBy      int
	products    []models.Product
	totalCnt    uint64
}

// NewGetProducts instance
func NewGetProducts(guid string, offset, limit int, searchQuery, category string, siteId string, featured string, sortBy int) *Products {
	ps := new(Products)
	ps.products = make([]models.Product, 0)
	ps.guid = guid
	ps.offset = offset
	ps.limit = limit
	ps.searchQuery = searchQuery
	ps.category = category
	ps.siteId = siteId
	ps.featured = featured
	ps.sortBy = sortBy
	return ps
}

// Data returns the final output from Do
func (ps *Products) Data() ([]models.Product, uint64) {
	return ps.products, ps.totalCnt
}

// Do tasks
func (ps *Products) Do() (err error) {

	t0 := time.Now()

	if err = ps.validate(); err != nil {
		return err
	}

	if err = ps.getProducts(); err != nil {
		return err
	}

	log.WithField("duration", time.Since(t0)).Info("admin product listing duration")

	t0 = time.Now()

	if err = ps.getProductsCount(); err != nil {
		return err
	}

	log.WithField("duration", time.Since(t0)).Info("admin product listing - determining products count")

	return nil
}

func (ps *Products) validate() (err error) {
	if ps.offset < 0 {
		err = errors.New("offset can not be negative")
		log.Warn(err.Error())
		return err
	}

	if ps.limit < 0 {
		err = errors.New("limit can not be negative")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (ps *Products) getProducts() (err error) {

	ps.products, err = models.GetProducts(
		ps.searchQuery,
		ps.offset,
		ps.limit,
		ps.category,
		ps.siteId,
		ps.guid,
		ps.sortBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ps *Products) getProductsCount() (err error) {
	ps.totalCnt, err = models.GetProductsCount(ps.searchQuery,
		ps.category,
		ps.siteId,
		ps.guid)
	if err != nil {
		return err
	}

	return nil
}
