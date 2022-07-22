package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ulventech/retro-ced-backend/models"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	form "github.com/thedanielforum/typed-form"

	services "github.com/ulventech/retro-ced-backend/services/admin"
	"github.com/ulventech/retro-ced-backend/services/product"
)

// Admin controller
type Admin struct {
	Base
}

// GetSites returns all sites that we are scraping
func (b *Admin) GetSites(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	s := services.NewGetSites()
	if err := s.Do(); err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}

// UpdateSite updates scrape site info
func (b *Admin) UpdateSite(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	// TODO: Don't do json
	var json services.Update
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	s := services.NewUpdate(json)
	if err = s.Do(); err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ClickTrackingSummary generates click tracker summary
func (b *Admin) ClickTrackingSummary(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, SITE_OWNER_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	s := services.NewClickTrackingSummary()
	if err := s.Do(); err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}

// GetProductsByApprovedStatus returns produts based on their approval status
func (b *Admin) GetProductsByApprovedStatus(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	// TODO Use daniel's awesome form lib (totaly not written by daniel...)
	itemsPerPage, err := strconv.ParseInt(c.Query("pageSize"), 10, 64)
	// var itemsPerPage int64 = 20
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	offset := (page - 1) * itemsPerPage
	searchQuery := c.Query("search")
	category := c.Query("category")
	siteId := c.Query("siteId")
	sortBy, _ := strconv.Atoi(c.Query("sortBy"))

	if siteId == "0" {
		siteId = ""
	}

	s := services.NewGetProductsByStatus(
		c.Param("status"),
		offset,
		itemsPerPage,
		searchQuery,
		category,
		siteId,
		sortBy,
	)

	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	items, count := s.Data()
	c.Header("x-total-count", fmt.Sprintf("%d", count))
	c.JSON(http.StatusOK, items)
}

func (b *Admin) UpdateApprovedStatus(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	var json services.UpdateApprovedStatus
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	up := services.NewUpdateApprovedStatus(json)
	if err = up.Do(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Admin) DeleteProductsById(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	var json services.DeleteProductsById
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	up := services.NewDeleteProductsById(json)
	if err = up.Do(); err != nil {

		log.WithError(err).
			WithField("Guid", up.ID).
			Info("could not delete product")

		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Admin) SetFeatured(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	var json product.SetFeaturedProducts
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}
	s := product.NewSetFeaturedProducts(json)
	if err = s.Do(); err != nil {
		c.JSON(http.StatusInternalServerError, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetProductByID returns product by product id
func (b *Admin) GetProductByID(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	// parse all the from data to types
	form := form.Parse(c)

	p, err := models.GetProductByGuid(form.GetParamString("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, p)
}

// UpdateProductById will update the specific product.
func (b *Admin) UpdateProductById(ctx *gin.Context) {

	user, err := b.Authenticate(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, b.errorMsg(ctx, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		ctx.JSON(http.StatusUnauthorized, b.errorMsg(ctx, "You are not authorized"))
		return
	}

	var req services.UpdateProductById
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	s := services.NewUpdateProductById(req)
	if err = s.Do(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "could not update product"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Admin) DeleteProductById(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, INTERN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	// parse all the from data to types
	form := form.Parse(c)

	s := services.NewDeleteProductById(
		form.GetParamInt64("product_id"),
	)
	if form.Errors() != nil {
		err := form.Errors()[0]
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Admin) GetUsers(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	offset := (page - 1) * pageSize
	term := c.Query("search")

	s := services.NewGetUsers(
		offset,
		pageSize,
		term,
	)
	if err = s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	users, totalCnt := s.Data()

	c.Header("x-total-count", fmt.Sprint(totalCnt))
	c.JSON(http.StatusOK, users)
}

// ChangeUserPermission changes users permissions
func (b *Admin) ChangeUserPermission(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	form := form.Parse(c)

	s := services.NewUpdateUser(
		form.GetParamInt64("userID"),
		form.GetInt64("acl"),
		form.GetInt64("site_id"),
	)
	if form.Errors() != nil {
		err := form.Errors()[0]
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	if err = s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}

// DeleteSite delete a site with all asisiated prodcuts and urls
func (b *Admin) DeleteSite(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, ADMIN_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	// parse all the from data to types
	form := form.Parse(c)

	s := services.NewDeleteSite(form.GetParamInt64("siteID"))
	if form.Errors() != nil {
		err := form.Errors()[0]
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	if err = s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetPartnerTrackingSummary returns a partners click tracking
func (b *Admin) GetPartnerTrackingSummary(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	s := services.NewGetPartnerTrackingSummary(user.UID)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}

// GetProductTrackingSummary returns a partners individual product click tracking
func (b *Admin) GetProductTrackingSummary(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	if !(b.Authorize(user.UID, SITE_OWNER_PERMISSION)) {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, "You are not authorized"))
		return
	}

	s := services.NewGetProductTrackingSummary(
		user.UID,
		c.Query("fromDate"),
		c.Query("toDate"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, s.Data())
}
