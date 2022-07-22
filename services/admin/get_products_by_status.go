package admin

import (
	"github.com/apex/log"

	"github.com/ulventech/retro-ced-backend/models"
)

// Products service
type Products struct {
	approved    bool
	searchQuery string
	category    string
	siteId      string
	offset      int64
	limit       int64
	sortBy      int
	products    []models.Product
	totalCnt    int64
}

// NewGetProductsByStatus instance
func NewGetProductsByStatus(status string, offset, limit int64, searchQuery, category, siteId string, sortBy int) *Products {
	n := new(Products)
	n.approved = status == "approved"
	n.offset = offset
	n.limit = limit
	n.searchQuery = searchQuery
	n.category = category
	n.siteId = siteId
	n.sortBy = sortBy
	return n
}

// Data return from Do function
func (p *Products) Data() ([]models.Product, int64) {
	return p.products, p.totalCnt
}

// Do tasks
func (p *Products) Do() (err error) {
	if err = p.getProductsByApprovedStatus(); err != nil {
		return err
	}

	if err = p.getProductsCountByStatus(); err != nil {
		return err
	}

	return nil
}

func (p *Products) getProductsByApprovedStatus() (err error) {
	p.products, err = models.GetProductsByApprovedStatus(
		p.approved,
		p.offset,
		p.limit,
		p.searchQuery,
		p.category,
		p.siteId,
		p.sortBy,
	)
	if err != nil {
		log.WithError(err).Error("failed to get products")
		return err
	}

	return nil
}

func (p *Products) getProductsCountByStatus() (err error) {
	p.totalCnt, err = models.GetProductsCountByStatus(p.approved, p.searchQuery, p.category, p.siteId)
	if err != nil {
		return err
	}

	return nil
}
