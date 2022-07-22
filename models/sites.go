package models

import (
	"time"

	"github.com/apex/log"
)

// Defines for known sites
const (
	SiteNameTradesy      = "tradesy.com"
	SiteNameMalleries    = "malleries.com"
	SiteNameAmoreVintage = "amorevintagejapan.com"
)

// Site model
type Site struct {
	ID          uint64    `json:"id" gorm:"column:id;primaryKey"`
	Name        string    `json:"name" gorm:"column:name"`
	URL         string    `json:"url" gorm:"column:url"`
	Sleep       int64     `json:"sleep" gorm:"column:sleep"`
	Active      bool      `json:"active" gorm:"column:active"`
	MaxPage     int64     `json:"max_page" gorm:"column:max_page"`
	LastScraped time.Time `json:"last_scraped" gorm:"column:last_scraped"`
}

// TableName returns the SQL table name for Site.
func (s *Site) TableName() string {
	return "Sites"
}

// GetByID returns site by ID
func (s *Site) GetByID(id int64) (*Site, error) {
	err := GetDB().QueryRow(`
		SELECT id, name, url, sleep, active, max_page, last_scraped
		FROM Sites
		WHERE id = ?
		`,
		id,
	).Scan(
		&s.ID,
		&s.Name,
		&s.URL,
		&s.Sleep,
		&s.Active,
		&s.MaxPage,
		&s.LastScraped,
	)
	if err != nil {
		log.WithError(err).Warn("get site by id failed")
		return s, err
	}

	return s, nil
}

// GetByName returns site by name
func (s *Site) GetByName(name string) (Site, error) {
	var site Site
	site.Name = name
	err := GetDB().QueryRow(`
		SELECT id, name, url, sleep, active, max_page, last_scraped
		FROM Sites
		WHERE name = ?
		ORDER BY id ASC
	`, name).Scan(
		&site.ID,
		&site.Name,
		&site.URL,
		&site.Sleep,
		&site.Active,
		&site.MaxPage,
		&site.LastScraped,
	)
	if err != nil {
		log.WithError(err).Warn("get site by name failed")
		return site, err
	}

	return site, nil
}

// GetAllSites returns the list of existing sites.
func GetAllSites() ([]Site, error) {
	var sites []Site
	err := GetDBv2().Find(&sites).Error
	if err != nil {
		return nil, err
	}

	return sites, nil
}

// Update sites
func (s *Site) Update(ID, maxPages uint64, sleep int64, name string, active bool) error {
	stmt, err := GetDB().Prepare(`
		ALTER TABLE Sites UPDATE name=?, sleep=?, active=?, max_page=?
		WHERE id = ?
	`)
	if err != nil {
		log.WithError(err).Error("failed prepare update user statement")
		return err
	}
	_, err = stmt.Exec(
		name,
		sleep,
		active,
		maxPages,
		ID,
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update user")
		return err
	}

	return nil
}

// UpdateLastScraped time
func (s *Site) UpdateLastScraped(ID uint64) error {
	stmt, err := GetDB().Prepare(`
		ALTER TABLE Sites UPDATE Sites
		SET last_scraped=?
		WHERE id=?
	`)
	if err != nil {
		log.WithError(err).Error("failed to prepare update last scraped statement")
		return err
	}
	_, err = stmt.Exec(
		time.Now().UTC(),
		ID,
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update last scraped")
		return err
	}

	return nil
}

func (s *Site) Delete(siteID int64) error {
	stmt, err := GetDB().Prepare(`
		ALTER TABLE Sites DELETE WHERE id = ?
	`)
	if err != nil {
		log.WithError(err).Error("failed prepare delete Site statement")
		return err
	}
	_, err = stmt.Exec(siteID)
	if err != nil {
		log.WithError(err).Error("failed to run exec on delete Site")
		return err
	}

	return nil
}
