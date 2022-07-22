package main

import (
	"encoding/json"
	"fmt"
	"github.com/ulventech/retro-ced-backend/models"

	"github.com/pkg/errors"
	"github.com/ulventech/retro-ced-backend/scraper/pages/malleries"
)

func fetchMalleries(address string) error {
	m, err := malleries.New(models.Site{
		ID: 10,
		MaxPage: 10,
		Sleep: 10,
	})
	if err != nil {
		return errors.Wrap(err, "could not create malleries crawler")
	}

	ad, err := m.Fetch(address)
	if err != nil {
		return errors.Wrap(err, "could not fetch malleries ad")
	}

	payload, _ := json.Marshal(ad)

	fmt.Printf("%s\n", payload)

	return nil
}
