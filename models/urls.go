package models

import (
	//"database/sql"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// Url model
type Url struct {
	Id          uint64 `db:"id"`
	SiteId      uint64 `db:"site_id"`
	Url         string `db:"url"`
	Category    string `db:"category"`
	LastUpdated time.Time `db:"last_updated"`
	CreatedAt   time.Time `db:"created_at"`
}

func (m *Url) TableName() string {
	return "Urls"
}

// AddUrl will store the given URL in the database
func (m *Url) AddUrl(siteId uint64, url string, category string) (err error) {
	if siteId <= 0 {
		return errors.New("siteId is required")
	}
	if category == "" {
		log.Fatalf("scraper url category is required: %s", category)
		return fmt.Errorf("category required")
	}

	// Check if there is a existing row
	//var u Url
	m.SiteId = siteId
	m.Url = url
	m.Category = category
	m.CreatedAt = time.Now()
	m.LastUpdated = time.Now()


	//err = GetDB().QueryRow(
	//	"SELECT id FROM urls WHERE site_id = ? AND url = ?",
	//	siteId,
	//	url,
	//).Scan(
	//	&u.Id,
	//)
	//if err != nil && err != sql.ErrNoRows {
	//	log.WithError(err).Error("fetch url failed")
	//	return err
	//}
	//m.Id = u.Id

	// Nothing exists so let's create a new one!
	//if err == sql.ErrNoRows {
	if err := GetDBv2().Create(m).Error; err != nil {
		log.Errorf("failed to insert url %s\n", err)
		return err
	}

		//stmt, err := GetDB().Prepare("INSERT INTO Urls(site_id, url, category, last_updated, created_at) VALUES(?, ?, ?, ?, ?)")
		//if err != nil {
		//	log.WithError(err).Error("failed prepare insert url statement")
		//	return err
		//}
		//item, err := stmt.Exec(siteId, url, category, time.Now(), time.Now())
		//if err != nil {
		//	log.WithError(err).Error("failed to run exec on insert url")
		//	return err
		//}
		//id, err := item.LastInsertId()
		//if err != nil {
		//	log.WithError(err).Error("failed to get last insert Guid")
		//}
		//m.Id = uint64(id)
		//
		//log.WithField("url", url).Debug("added new url")

	//} else { // Let's set the last_updated so we know it was checked
	//	stmt, err := GetDB().Prepare("UPDATE Urls SET last_updated = ?, category = ? WHERE id = ?")
	//	if err != nil {
	//		log.WithError(err).Error("failed prepare update url statement")
	//		return err
	//	}
	//	_, err = stmt.Exec(time.Now(), category, u.Id)
	//	if err != nil {
	//		log.WithError(err).Error("failed to run exec on update url")
	//		return err
	//	}
	//
	//	log.WithField("url", url).Debug("updated existing url")
	//
	//}

	return nil
}

func (m *Url) GetUrls(siteId uint64, category string) (urls []Url, err error) {
	rows, err := GetDB().Query("SELECT id, url FROM Urls WHERE site_id = ? AND category = ? ORDER BY created_at DESC", siteId, category)
	if err != nil {
		log.WithError(err).Error("failed to fetch urls")
		return urls, err
	}

	defer rows.Close()
	for rows.Next() {
		var u Url
		err := rows.Scan(&u.Id, &u.Url)
		if err != nil {
			log.WithError(err).Error("failed to scan url")
			return urls, err
		}
		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		log.WithError(err).Error("rows error for url")
		return urls, err
	}

	return urls, nil
}

func (m *Url) UpdateUrl(url_id int64, url string) (err error) {
	if url_id <= 0 {
		return errors.New("UrlId is required")
	}

	stmt, err := GetDB().Prepare(`UPDATE Urls
		SET url = ?
		WHERE id = ?`)
	if err != nil {
		log.WithError(err).Error("failed prepare update url statement")
		return err
	}
	_, err = stmt.Exec(url, url_id)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update url")
		return err
	}
	return nil
}
