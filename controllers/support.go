package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	form "github.com/thedanielforum/typed-form"

	"github.com/ulventech/retro-ced-backend/services/services"
)

// Support controller
type Support struct {
	Base
}

// ContactUs sends a email to site admin with form data
func (b *Support) ContactUs(c *gin.Context) {
	form := form.Parse(c)
	s := services.NewContactUs(
		form.GetString("email"),
		form.GetString("name"),
		form.GetString("phone"),
		form.GetString("message"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// PartnerApplication sends a email to site admin with form data
func (b *Support) PartnerApplication(c *gin.Context) {
	form := form.Parse(c)
	s := services.NewPartnerApplication(
		form.GetString("platform"),
		form.GetString("otherPlatform"),
		form.GetString("firstName"),
		form.GetString("lastName"),
		form.GetString("email"),
		form.GetString("storeName"),
		form.GetString("website"),
		form.GetString("location"),
		form.GetString("inventorySystem"),
		form.GetString("numberOfItems"),
		form.GetString("topFiveBrands"),
		form.GetString("affiliateNetwork"),
		form.GetString("otherAffiliateNetwork"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
