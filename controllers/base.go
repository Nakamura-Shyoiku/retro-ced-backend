package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/ulventech/retro-ced-backend/services/auth"
	"github.com/ulventech/retro-ced-backend/utils"
)

// Base is the root controller
type Base struct {
	Trace string
}

const (
	SHOPPER_PERMISSION    = 0
	INTERN_PERMISSION     = 10
	SITE_OWNER_PERMISSION = 100
	ADMIN_PERMISSION      = 1000
)

// Error struct
type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Trace   string `json:"trace"`
}

func (b *Base) errorMsg(c *gin.Context, msg string) gin.H {
	return gin.H{
		"error": &Error{
			Message: msg,
			Trace:   c.GetString("trace"),
		},
	}
}

// Authenticate checks JWT token's validity
func (b *Base) Authenticate(c *gin.Context) (*utils.JwtCustomClaims, error) {
	// Get JWT token from header
	token := c.GetHeader("Authorization")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return claims, err
	}

	return claims, nil
}

// Authorize checks JWT token's permission
func (b *Base) Authorize(userID string, permissionID int64) bool {
	// Get JWT token from header
	return auth.HasPermission(userID, permissionID)
}
