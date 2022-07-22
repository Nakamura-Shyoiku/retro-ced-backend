package controllers

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/query"
	"github.com/ulventech/retro-ced-backend/services/product"
)

type Product struct {
	Base
}

func (b *Product) GetProduct(c *gin.Context) {
	p := product.NewGetProduct(c.Param("id"))
	if err := p.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, p.Data())
}

func (b *Product) GetProducts(c *gin.Context) {

	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	var params struct {
		Guid        string `form:"guid"`
		Page        int    `form:"page"`
		PageSize    int    `form:"pageSize"`
		SearchQuery string `form:"search"`
		Category    string `form:"category"`
		Offset      int    `form:"offset"`
		SiteID      string `form:"siteId"`
		Featured    string `form:"featured"`
		SortBy      int    `form:"sortBy"`
	}

	err = c.BindQuery(&params)
	if err != nil {
		log.WithError(err).Error("invalid search params")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	offset := (params.Page - 1) * params.PageSize

	p := product.NewGetProducts(
		params.Guid,
		offset,
		params.PageSize,
		params.SearchQuery,
		params.Category,
		params.SiteID,
		params.Featured,
		params.SortBy,
	)

	if err := p.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	items, count := p.Data()
	convCount := strconv.FormatUint(count, 10)
	c.Header("x-total-count", convCount)
	c.JSON(http.StatusOK, items)
}

// GetProductsCount returns the number of products available for /admin/products list, for a given search.
func (b *Product) GetProductsCount(ctx *gin.Context) {

	// TODO: this should be protected - admin only

	var params struct {
		SearchQuery string `form:"search"`
		Category    string `form:"category"`
		SiteID      string `form:"siteId"`
		Featured    string `form:"featured"`
	}

	err := ctx.BindQuery(&params)
	if err != nil {
		return
	}

	count, err := models.GetProductsCount(params.SearchQuery, params.Category, params.SiteID, params.Featured)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "could not retrieve product count"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

// GetProductsByFeatured will return featured products of a certain category.
func (b *Product) GetProductsByFeatured(ctx *gin.Context) {

	featuredCategory := ctx.Query(query.FeaturedParamName)
	if featuredCategory == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fetchSize, err := strconv.Atoi(ctx.Query("fetchSize"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO: what should we do here
	// user, _ := b.Authenticate(ctx)

	// userRecord, err := models.GetUserByGUID(user.UID)
	// if err != nil {
	// 	ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "could not find user"))
	// 	return
	// }

	products, err := models.GetProductsByFeatured(featuredCategory, fetchSize)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "could not lookup favourites"))
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (b *Product) GetProductsByCategory(c *gin.Context) {
	user, _ := b.Authenticate(c)
	p := product.NewGetProductsByCategory(user.UID)
	if err := p.Do(
		c.Param("category"),
		c.Param("offset"),
		c.QueryArray("category"),
		c.QueryArray("brand"),
		c.QueryArray("color"),
		c.QueryArray("size"),
		c.QueryArray("shoeSize")); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	items, count := p.Data()
	convCount := strconv.FormatUint(count, 10)
	c.Header("x-total-count", convCount)
	c.JSON(http.StatusOK, items)
}

func (b *Product) GetProductsByBrand(c *gin.Context) {
	user, _ := b.Authenticate(c)

	p := product.NewGetProductsByBrand(user.UID)
	if err := p.Do(c.Param("brand"), c.Param("offset")); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	items, count := p.Data()
	c.Header("x-total-count", fmt.Sprintf("%d", count))
	c.JSON(http.StatusOK, items)
}

func (b *Product) GetBrands(c *gin.Context) {
	p := product.NewBrands()
	if err := p.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	c.JSON(http.StatusOK, p.Data())
}

func (b *Product) GetFavourites(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	p := product.NewFavourites(user.UID)
	if err := p.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	items, count := p.Data()
	convCount := strconv.FormatUint(count, 10)
	c.Header("x-total-count", convCount)
	c.JSON(http.StatusOK, items)
}

func (b *Product) AddFavourite(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	p := product.NewAddFavourite(user.UID)
	i, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err = p.Do(i); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Product) RemoveFavourite(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	i, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	p := product.NewRemoveFavourite(user.UID)
	if err = p.Do(i); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Product) IsFavourited(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	pid, _ := strconv.ParseInt(c.Param("product_id"), 10, 64)

	s := product.NewIsFavourite(pid, user.UID)
	if err = s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favourited": s.Data()})
}

func (b *Product) SearchProducts(c *gin.Context) {
	user, _ := b.Authenticate(c)

	p := product.NewSearchProducts(user.UID)
	if err := p.Do(c.Param("search"), c.Param("offset")); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	items, count := p.Data()
	convCount := strconv.FormatUint(count, 10)
	c.Header("x-total-count", convCount)
	c.JSON(http.StatusOK, items)
}
