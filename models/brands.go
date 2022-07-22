package models

import (
	"database/sql"

	"github.com/apex/log"
)

type Brand struct {
	Brand string `json:"brand"`
}

func GetBrands() ([]Brand, error) {
	var bs []Brand
	rows, err := GetDB().Query(`SELECT DISTINCT brand FROM Products ORDER BY brand`)
	if err != nil {
		log.WithError(err).Error("failed to query brands")
	}
	defer rows.Close()
	for rows.Next() {
		var b Brand
		err := rows.Scan(
			&b.Brand,
		)
		if err != nil {
			log.WithError(err).Error("failed to get brands by rows")
			return bs, err
		}
		bs = append(bs, b)
	}
	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("brands row error")
		return bs, err
	}
	return bs, nil
}
