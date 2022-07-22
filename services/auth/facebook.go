package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/spf13/viper"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/utils"
)

// Facebook service
type Facebook struct {
	code        string
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	FbID        string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`

	user models.User
}

// NewFacebook instance
func NewFacebook(code string) *Facebook {
	n := new(Facebook)
	n.code = code
	return n
}

// GetToken returns token from Do function
func (f *Facebook) GetToken() string {
	return utils.NewToken(
		f.user.Guid,
		f.Email,
		"",
		f.user.FirstName,
		f.user.LastName,
		f.user.FbID,
		f.user.ACL,
	)
}

// Do tasks
func (f *Facebook) Do() (err error) {
	if err = f.exchangeCode(); err != nil {
		return err
	}

	if err = f.getFbUser(); err != nil {
		return err
	}

	if err = f.createUser(); err != nil {
		return err
	}

	return nil
}

func (f *Facebook) exchangeCode() error {
	// TODO Use custom client, so we set a timeout
	resp, err := http.Get(fmt.Sprintf(
		"https://graph.facebook.com/v2.11/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		viper.GetString("facebook.app_id"),
		viper.GetString("facebook.app_secret"),
		f.code,
		fmt.Sprintf("%s/v1/auth/fb/callback", viper.GetString("app.domain")),
	))
	if err != nil {
		log.WithError(err).Error("failed to exchange code for access token")
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(f)
	if err != nil {
		log.WithError(err).Error("failed to decode response")
		return err
	}

	return nil
}

func (f *Facebook) getFbUser() error {
	client := &http.Client{Timeout: time.Duration(3600 * time.Second)}
	req, _ := http.NewRequest("GET", "https://graph.facebook.com/v2.11/me?fields=email,first_name,last_name", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", f.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to get user info")
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(f)
	if err != nil {
		log.WithError(err).Error("failed to decode response")
		return err
	}

	return nil
}

func (f *Facebook) createUser() (err error) {
	f.user, err = models.GetUserByEmail(f.Email)
	if err != nil {
		log.WithError(err).Warn("error while fetching user")
		return err
	}
	// User don't exist, make new one
	if f.user.Id <= 0 {
		log.Info("creating new user")
		u := models.User{
			FirstName: f.FirstName,
			LastName:  f.LastName,
			Email:     f.Email,
			FbID:      f.FbID,
		}
		f.user.Guid, err = u.Create()
		if err != nil {
			return err
		}
		// Subscribe the user
		//if err = f.mailChimp(); err != nil {
		//	log.WithError(err).Error("error while subscribing")
		//}
	} else {
		// Set the facebook Guid
		if err = models.UpdateFbID(f.user.Guid, f.FbID); err != nil {
			return err
		}
	}

	return nil
}

func (f *Facebook) mailChimp() (err error) {
	err = mailchimp.SetKey(viper.GetString("mailchimp.key"))
	if err != nil {
		return err
	}

	params := &members.NewParams{
		EmailAddress: f.Email,
		Status:       members.StatusSubscribed,
	}

	_, err = members.New(viper.GetString("mailchimp.members"), params)
	if err != nil {
		log.WithError(err).Error("failed to add email to list")
		return err
	}

	return nil
}
