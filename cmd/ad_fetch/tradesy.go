package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/ulventech/retro-ced-backend/models"
	"github.com/ulventech/retro-ced-backend/scraper/pages/tradesy"
)

func fetchTradesy(address string) error {

	cfg, err := new(models.Site).GetByName("tradesy.com")
	if err != nil {
		return errors.Wrap(err, "could not get site config")
	}

	t, err := tradesy.New(cfg)
	if err != nil {
		return errors.Wrap(err, "could not create tradesy crawler")
	}

	payload, err := t.Fetch(address)
	if err != nil {
		return errors.Wrap(err, "could not fetch tradsy ad")
	}

	fmt.Printf("%s\n", payload)

	return nil
}
