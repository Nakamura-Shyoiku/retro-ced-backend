package models

import (
	"time"

	"github.com/apex/log"
)

// Click is data that contains information on the pages that was clicked on the site that will be referred to affliate sites
type Click struct {
	ID        uint64    `json:"id"`
	Link      string    `json:"link"`
	Timestamp time.Time `json:"timestamp"`
	Clicks    int64     `json:"clicks"`
}

func (c *Click) TableName() string {
	return "Clicks"
}

// Create command to save a click value into Database
func (c *Click) Create(link string) error {
	stmt, err := GetDB().Prepare(`INSERT INTO Clicks(link) VALUES (?)`)
	if err != nil {
		log.WithError(err).Error("failed to prepare Click Create")
		return err
	}

	_, err = stmt.Exec(link)
	if err != nil {
		log.WithError(err).Error("failed to execute Click Create")
		return err
	}

	return nil
}
