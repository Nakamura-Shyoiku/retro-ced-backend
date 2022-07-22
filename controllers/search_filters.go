package controllers

import (
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
)

var badSearchParamsError = errors.Errorf("invalid search parameter")

// SearchBrand will return available brands for the given search request
func (b *Search) SearchBrand(c *gin.Context) {

	fields, err := b.getDistinctFields(c, "brand")
	if err != nil {

		log.WithError(err).Warn("could not get available brands")

		code := http.StatusInternalServerError
		if errors.Cause(err) == badSearchParamsError {
			code = http.StatusBadRequest
		}

		c.JSON(code, b.errorMsg(c, "could not retrieve available values"))
		return
	}

	c.JSON(http.StatusOK, fields)
}

// SearchSize will return available sizes for the given search request
func (b *Search) SearchSize(c *gin.Context) {

	fields, err := b.getDistinctFields(c, "size")
	if err != nil {

		log.WithError(err).Warn("could not get available sizes")

		code := http.StatusInternalServerError
		if errors.Cause(err) == badSearchParamsError {
			code = http.StatusBadRequest
		}

		c.JSON(code, b.errorMsg(c, "could not retrieve available values"))
		return
	}

	c.JSON(http.StatusOK, fields)
}

// SearchSubcategory will return available subcategories for the given search request
func (b *Search) SearchSubcategory(c *gin.Context) {

	fields, err := b.getDistinctFields(c, "sub_category")
	if err != nil {

		log.WithError(err).Warn("could not get available subcategories")

		code := http.StatusInternalServerError
		if errors.Cause(err) == badSearchParamsError {
			code = http.StatusBadRequest
		}

		c.JSON(code, b.errorMsg(c, "could not retrieve available values"))
		return
	}

	c.JSON(http.StatusOK, fields)
}

// SearchColor will return available colors for the given search request
func (b *Search) SearchColor(c *gin.Context) {

	fields, err := b.getDistinctFields(c, "color")
	if err != nil {

		log.WithError(err).Warn("could not get available colors")

		code := http.StatusInternalServerError
		if errors.Cause(err) == badSearchParamsError {
			code = http.StatusBadRequest
		}

		c.JSON(code, b.errorMsg(c, "could not retrieve available values"))
		return
	}

	c.JSON(http.StatusOK, fields)
}

func (b *Search) getDistinctFields(c *gin.Context, column string) ([]string, error) {

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
			return nil, badSearchParamsError
		}

		params.priceMin = uint(min)
	}

	if priceMax != "" {
		max, err := strconv.ParseUint(priceMax, 10, 32)
		if err != nil {
			return nil, badSearchParamsError
		}

		params.priceMax = uint(max)
	}

	query, args := params.getUniqueFieldsQuery(column)
	log.WithField("query", query).
		WithField("args", args).
		WithField("filter", column).
		Info("built query for available filters")

	rows, err := models.GetDB().Query(query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "could not execute query (query: %v| args: %v)", query, args)
	}

	defer rows.Close()

	fields := make([]string, 0)
	for rows.Next() {

		var field string
		err = rows.Scan(&field)
		if err != nil {
			log.WithError(err).
				WithField("filter", column).
				Warn("could not scan filter value")

			continue
		}

		if field == "" {
			continue
		}

		fields = append(fields, field)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "error iterating search results (query: %v, args: %v)", query, args)
	}

	return fields, nil
}

func (sp *searchParams) getUniqueFieldsQuery(column string) (string, []interface{}) {

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select(column)
	sb.Distinct()
	sb.From("Products")

	// exact matches
	if sp.category != "" {
		sb.Where(sb.Equal("Products.category", sp.category))
	}

	if len(sp.colors) > 0 {
		sb.Where(sb.In("Products.color", sqlbuilder.Flatten(sp.colors)...))
	}

	if len(sp.sizes) > 0 {
		sb.Where(sb.In("Products.size", sqlbuilder.Flatten(sp.sizes)...))
	}

	if len(sp.subcategories) > 0 {
		sb.Where(sb.In("Products.sub_category", sqlbuilder.Flatten(sp.subcategories)...))
	}

	if len(sp.brands) > 0 {
		sb.Where(sb.In("Products.brand", sqlbuilder.Flatten(sp.brands)...))
	}

	// ranges
	if sp.priceMin != 0 {
		sb.Where(sb.GreaterEqualThan("Products.price", sp.priceMin))
	}

	if sp.priceMax != 0 {
		sb.Where(sb.LessEqualThan("Products.price", sp.priceMax))
	}

	// like - search for "%searchterm%"
	if sp.queryText != "" {
		sb.Where(sb.Like("Products.title", `%`+sp.queryText+`%`))
	}

	return sb.Build()
}
