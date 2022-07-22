package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

const (
	querySearchParam       = "query"
	categorySearchParam    = "category"
	colorSearchParam       = "color"
	priceMinSearchParam    = "pricemin"
	priceMaxSearchParam    = "pricemax"
	sizeSearchParam        = "size"
	subcategorySearchParam = "sub_category"
	pageSearchParam        = "page"
	countSearchParam       = "count"
	brandSearchParam       = "brand"

	maxHitsPerPage = 100
)

// searchResult is a models.Product - slightly reformatted and with certain fields omitted
type searchResult struct {
	Guid         string `json:"guid"`
	Brand        string `json:"brand"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	URL          string `json:"url"`
	SubCategory  string `json:"sub_category"`
	Color        string `json:"color"`
	Image        string `json:"image"`
	Model        string `json:"model"`
	Price        int64  `json:"price"`
	Condition    string `json:"condition"`
	ShoeSize     string `json:"shoe_size"`
	Size         string `json:"size"`
	Measurements string `json:"measurements"`
	RetailPrice  int64  `json:"retail_price"`
}

// Search controller
type Search struct {
	Base
}

// get multiple fields, separated by commas; filter out empty values
func getMultipleFilters(value string) []string {

	fields := strings.Split(value, ",")

	out := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			out = append(out, strings.ToLower(field))
		}
	}

	return out
}

// SearchProducts will perform a DB search using the given parameters
func (b *Search) SearchProducts(c *gin.Context) {

	// supported query fields:
	// 1. query (partial match on title)
	// 2. category
	// 3. color
	// 4. price (min, max)
	// 5. size
	// 6. subcategory (e.g. bags can be clutches)
	// 7. brand

	// other parameters:
	// 7. count - how many items on page, max 100
	// 8. page - which page to show

	params := searchParams{
		queryText:     c.PostForm(querySearchParam),
		category:      c.PostForm(categorySearchParam),
		subcategories: getMultipleFilters(c.PostForm(subcategorySearchParam)),
		colors:        getMultipleFilters(c.PostForm(colorSearchParam)),
		sizes:         getMultipleFilters(c.PostForm(sizeSearchParam)),
		brands:        getMultipleFilters(c.PostForm(brandSearchParam)),
	}

	// get price range
	priceMin := c.PostForm(priceMinSearchParam)
	priceMax := c.PostForm(priceMaxSearchParam)

	if priceMin != "" {
		min, err := strconv.ParseUint(priceMin, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, b.errorMsg(c, "invalid price"))
			return
		}

		params.priceMin = uint(min)
	}

	if priceMax != "" {
		max, err := strconv.ParseUint(priceMax, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, b.errorMsg(c, "invalid price"))
			return
		}

		params.priceMax = uint(max)
	}

	countParam := c.PostForm(countSearchParam)
	pageParam := c.PostForm(pageSearchParam)

	// set number of items to return
	if countParam != "" {
		n, err := strconv.ParseUint(countParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, b.errorMsg(c, "invalid count"))
			return
		}

		params.count = uint(n)
	}

	// if the number of items is not specified or is too large, use max
	if params.count == 0 || params.count > maxHitsPerPage {
		params.count = maxHitsPerPage
	}

	// set the number of page to return
	if pageParam != "" {
		n, err := strconv.ParseUint(pageParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, b.errorMsg(c, "invalid page"))
			return
		}

		params.page = uint(n)
	}

	products, err := params.search()
	if err != nil {
		log.WithError(err).
			WithField("search", params.String()).
			Warn("search failed")

		c.JSON(http.StatusInternalServerError, b.errorMsg(c, "search failed"))
		return
	}

	var res []searchResult
	for _, p := range products {
		var r = searchResult{
			p.Guid,
			p.Brand,
			p.Title,
			p.Description,
			p.Category,
			p.ProductURL,
			p.SubCategory,
			p.Color,
			p.Img,
			p.Model,
			p.Price,
			p.ProductCondition,
			p.ShoeSize,
			p.Size,
			p.Measurements,
			p.RetailPrice}
		res = append(res, r)
	}

	log.Infof("returning %d products from mysql", len(res))
	log.Infof("find products for category %s ", params.category)

	c.JSON(http.StatusOK, res)
}
