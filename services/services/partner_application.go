package services

import (
	"errors"
	
	"github.com/apex/log"
	
	"github.com/ulventech/retro-ced-backend/email"
)

// ContactUs service
type PartnerApplication struct {
	platform              string
	otherPlatform         string
	firstName             string
	lastName              string
	email                 string
	storeName             string
	website               string
	location              string
	inventorySystem       string
	numberOfItems         string
	topFiveBrands         string
	affiliateNetwork      string
	otherAffiliateNetwork string
}

// NewPartnerApplication instance
func NewPartnerApplication(platform, otherPlatform, firstName, lastName, email, storeName, website, location, inventorySystem, numberOfItems, topFiveBrands, affiliateNetwork, otherAffiliateNetwork string) *PartnerApplication {
	return &PartnerApplication{
		platform,
		otherPlatform,
		firstName,
		lastName,
		email,
		storeName,
		website,
		location,
		inventorySystem,
		numberOfItems,
		topFiveBrands,
		affiliateNetwork,
		otherAffiliateNetwork,
	}
}

// Do tasks
func (pa *PartnerApplication) Do() (err error) {
	if err = pa.validate(); err != nil {
		return err
	}

	if err = pa.sendMessage(); err != nil {
		return err
	}

	return nil
}

func (pa *PartnerApplication) validate() (err error) {
	if pa.platform == "" {
		err = errors.New("Platform is required")
		log.Warn(err.Error())
		return err
	}

	if pa.firstName == "" {
		err = errors.New("First name is required")
		log.Warn(err.Error())
		return err
	}

	if pa.lastName == "" {
		err = errors.New("Last name is required")
		log.Warn(err.Error())
		return err
	}

	if pa.email == "" {
		err = errors.New("Email is required")
		log.Warn(err.Error())
		return err
	}

	if pa.storeName == "" {
		err = errors.New("Store name is required")
		log.Warn(err.Error())
		return err
	}
	if pa.website == "" {
		err = errors.New("Website is required")
		log.Warn(err.Error())
		return err
	}

	if pa.location == "" {
		err = errors.New("Store location is required")
		log.Warn(err.Error())
		return err
	}

	if pa.inventorySystem == "" {
		err = errors.New("Inventory system is required")
		log.Warn(err.Error())
		return err
	}

	if pa.numberOfItems == "" {
		err = errors.New("Number of items is required")
		log.Warn(err.Error())
		return err
	}

	if pa.topFiveBrands == "" {
		err = errors.New("Brand is required")
		log.Warn(err.Error())
		return err
	}

	if pa.affiliateNetwork == "" {
		err = errors.New("Affiliate network is required")
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (pa *PartnerApplication) sendMessage() error {
	err := email.NewEmail().SendPartnerApplication(pa.email, email.Data{
		"Platform":              pa.platform,
		"OtherPlatform":         pa.otherPlatform,
		"FirstName":             pa.firstName,
		"LastName":              pa.lastName,
		"Email":                 pa.email,
		"StoreName":             pa.storeName,
		"Website":               pa.website,
		"Location":              pa.location,
		"InventorySystem":       pa.inventorySystem,
		"NumberofItems":         pa.numberOfItems,
		"TopFiveBrands":         pa.topFiveBrands,
		"AffiliateNetwork":      pa.affiliateNetwork,
		"OtherAffiliateNetwork": pa.otherAffiliateNetwork,
	})
	if err != nil {
		log.WithError(err).Error("failed to send partner application")
		return err
	}

	return nil
}
