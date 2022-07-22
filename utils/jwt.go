package utils

import (
	"time"

	"github.com/apex/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/models"
)

// JwtCustomClaims Struct that holds information on data that is passed to front end when it comes to user information
type JwtCustomClaims struct {
	UID       string `json:"uid"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ACL       int    `json:"acl"`
	FbID      string `json:"fbID"`
	jwt.StandardClaims
}

// NewToken Creates a new jwt token with the given parameters
func NewToken(uid, email, username, firstName, lastName, fbid string, acl int) string {
	exp := viper.GetInt("app.jwt_expire")
	claims := JwtCustomClaims{
		uid,
		email,
		username,
		firstName,
		lastName,
		acl,
		fbid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(exp) * time.Second).Unix(),
			Issuer:    "retroced",
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("app.jwt_secret")))
	if err != nil {
		log.WithError(err).Error("error while signing the JWT")
		return ""
	}

	return tokenString
}

// ValidateToken Validates the JWT token
func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	// sample token is expired.  override time so it parses as valid
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("app.jwt_secret")), nil
	})
	if err != nil {
		log.WithError(err).Error("jwt error")
		return new(JwtCustomClaims), err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	user, err := models.GetUserByEmail(claims.Email)
	if err != nil {
		log.WithError(err).Error("Unable to fetch error values from database")
		return new(JwtCustomClaims), err
	}

	if ok && token.Valid {
		return &JwtCustomClaims{
			UID:       user.Guid,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ACL:       user.ACL,
			FbID:      user.FbID,
		}, nil
	}

	return new(JwtCustomClaims), errors.New("jwt error")
}
