package controllers

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
)

type searchParams struct {
	queryText string
	category  string
	priceMin  uint
	priceMax  uint

	colors        []string
	sizes         []string
	subcategories []string
	brands        []string

	// pagination
	count uint
	page  uint
}

// for debugging purposes only
func (sp *searchParams) String() string {
	return fmt.Sprintf("<queryText: %v|category: %v|colors: %+v|priceMin: %v|priceMax: %v|sizes: %+v|subcategories: %+v|brands: %+v|page: %v|count: %v>",
		sp.queryText,
		sp.category,
		sp.colors,
		sp.priceMin,
		sp.priceMax,
		sp.sizes,
		sp.subcategories,
		sp.brands,
		sp.page,
		sp.count,
	)
}

func (sp *searchParams) search() ([]models.Product, error) {
	db := models.GetDBv2().Table("Products").Limit(int(sp.count))

	var offset int
	// translate page number to numbering that starts at zero
	if sp.page > 1 {
		offset = int((sp.page - 1) * sp.count)
	}

	db = db.Offset(offset)
	// exact matches
	if sp.category != "" {
		db = db.Where("lowerUTF8(category) = ?", sp.category)
	}

	if len(sp.colors) > 0 {
		db = db.Where("lowerUTF8(color) IN ?", sp.colors)
	}

	if len(sp.sizes) > 0 {
		db = db.Where("lowerUTF8(size) IN ?", sp.sizes)
	}

	if len(sp.subcategories) > 0 {
		db = db.Where("lowerUTF8(sub_category) IN ?", sp.subcategories)
	}

	if len(sp.brands) > 0 {
		db = db.Where("lowerUTF8(brand) IN ?", sp.brands)
	}

	// ranges
	if sp.priceMin != 0 {
		db = db.Where("price >= ?", sp.priceMin)
	}

	if sp.priceMax != 0 {
		db = db.Where("price <= ?", sp.priceMax)
	}

	// like - search for "%searchterm%"
	if sp.queryText != "" {
		db = db.Where("title ilike ?", `%`+sp.queryText+`%`)
	}
	db = db.Order("created_at desc")

	var prods []models.Product

	if err := db.Scan(&prods).Error; err != nil {
		return nil, errors.Wrap(err, "error querying products")
	}

	return prods, nil
}
