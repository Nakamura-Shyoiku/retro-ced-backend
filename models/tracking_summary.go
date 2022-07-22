package models

import (
	"github.com/apex/log"
)

// Tracking struct
type Tracking struct {
	Link                string `json:"link"`
	Clicks              int    `json:"clicks"`
	TodayClicks         int    `json:"today_clicks"`
	LastSevenDaysClicks int    `json:"last_seven_days_clicks"`
	ThisMonthClicks     int    `json:"this_month_clicks"`
	LastMonthClicks     int    `json:"last_month_clicks"`
	LastTwoMonthsClicks int    `json:"before_last_month_clicks"`
}

// GetTrackingSummary get clicks summary for partner sites
func (t *Tracking) GetTrackingSummary(siteName string) (Tracking, error) {
	siteNamehQuery := "%" + siteName + "%"
	var tracking Tracking
	err := GetDB().QueryRow(`
		SELECT SUM(DATE(Clicks.timestamp) = CURDATE()) AS today,
	   	SUM(DATE(Clicks.timestamp) >= CURDATE() - INTERVAL 6 DAY) as last_7_days,
       	SUM(YEAR(Clicks.timestamp) = YEAR(CURRENT_DATE) AND MONTH(Clicks.timestamp) = MONTH(CURRENT_DATE)) as this_month,
       	SUM(YEAR(Clicks.timestamp) = YEAR(CURRENT_DATE - INTERVAL 1 MONTH) AND MONTH(Clicks.timestamp) = MONTH(CURRENT_DATE - INTERVAL 1 MONTH))  as last_month,
       	SUM(YEAR(Clicks.timestamp) = YEAR(CURRENT_DATE - INTERVAL 2 MONTH) AND MONTH(Clicks.timestamp) = MONTH(CURRENT_DATE - INTERVAL 2 MONTH)) as before_last_month
		FROM Clicks
		WHERE link LIKE ?
	`, siteNamehQuery).Scan(
		&tracking.TodayClicks,
		&tracking.LastSevenDaysClicks,
		&tracking.ThisMonthClicks,
		&tracking.LastMonthClicks,
		&tracking.LastTwoMonthsClicks,
	)
	if err != nil {
		log.WithError(err).Error("fetch summary failed")
		return tracking, err
	}

	return tracking, nil
}

// GetTrackingSummary get clicks summary for partner sites individual product
func (t *Tracking) GetProductTrackingSummary(siteName, fromDate, toDate string) ([]Tracking, error) {
	siteNamehQuery := "%" + siteName + "%"
	var tracking []Tracking
	rows, err := GetDB().Query(`
		SELECT link, COUNT(*) AS num
		FROM Clicks
		WHERE link LIKE ?
		AND (DATE(timestamp) >= ? OR ? = '') AND (DATE(timestamp) <= ? OR ? = '')
		GROUP BY link
		ORDER BY num DESC
	`, siteNamehQuery,
		fromDate,
		fromDate,
		toDate,
		toDate,
	)
	if err != nil {
		log.WithError(err).Error("failed to get product tracking summary")
		return tracking, err
	}
	defer rows.Close()
	for rows.Next() {
		var t Tracking
		err := rows.Scan(
			&t.Link,
			&t.Clicks,
		)
		if err != nil {
			log.WithError(err).Error("failed to get product tracking summary rows")
			return tracking, err
		}
		tracking = append(tracking, t)
	}
	err = rows.Err()
	if err != nil {
		log.WithError(err).Error("get product tracking summary rows error")
		return tracking, err
	}

	return tracking, nil
}
