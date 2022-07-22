package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/services/auth"
	"github.com/ulventech/retro-ced-backend/utils"
)

// Auth controller
type Auth struct {
	Base
}

func (b *Auth) FbCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Warn("facebook callback code is required")
		c.JSON(http.StatusBadRequest, b.errorMsg(c, "code is required"))
		return
	}

	s := auth.NewFacebook(code)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	// switch between mobile and desktop site
	var rdr string
	if strings.Contains(c.Query("state"), "mobile_") {
		rdr = fmt.Sprintf(
			"http://%s/?token=%s",
			viper.GetString("auth.redirect_domain_mobile"),
			s.GetToken(),
		)
	} else {
		rdr = fmt.Sprintf(
			"http://%s/?token=%s",
			viper.GetString("auth.redirect_domain"),
			s.GetToken(),
		)
	}

	c.Redirect(http.StatusPermanentRedirect, rdr)
}

// ValidateToken validates token
func (b *Auth) ValidateToken(c *gin.Context) {
	user, err := b.Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": utils.NewToken(
			user.UID,
			user.Email,
			user.Username,
			user.FirstName,
			user.LastName,
			user.FbID,
			user.ACL,
		),
	})
}

// Login using email and password
func (b *Auth) Login(c *gin.Context) {
	s := auth.NewLogin(
		c.PostForm("email"),
		c.PostForm("password"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": s.Data()})
}

// Register using email and password
func (b *Auth) Register(c *gin.Context) {
	s := auth.NewRegister(
		c.PostForm("email"),
		c.PostForm("password"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": s.Data()})
}

// GetUserDetails returns users details
func (b *Auth) GetUserDetails(c *gin.Context) {
	// auth is optional for this endpoint
	user, _ := b.Authenticate(c)

	u := auth.NewGetUser(user.UID)
	if err := u.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, u.Data())
}

// ResetPassword sends user a password reset email
func (b *Auth) ResetPassword(c *gin.Context) {
	s := auth.NewPasswordReset(c.PostForm("email"))
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (b *Auth) SetNewPassword(c *gin.Context) {
	s := auth.NewSetPassword(
		c.PostForm("email"),
		c.PostForm("token"),
		c.PostForm("password"),
	)
	if err := s.Do(); err != nil {
		c.JSON(http.StatusBadRequest, b.errorMsg(c, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
