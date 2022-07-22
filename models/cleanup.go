package models

import (
	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	// removeStaleProducts = `DELETE FROM Products
	// 						WHERE id <> 0 AND id IN
	// 							(
	// 								SELECT stale.id
	// 						 	 	FROM (
	// 							 		SELECT id
	// 							 		FROM Products
	// 							 		WHERE last_updated < DATE_SUB( CURRENT_TIMESTAMP(), INTERVAL 6 MONTH)
	// 								) AS stale
	// 							)`

	// removeStaleURLs = `DELETE FROM Urls
	// 					WHERE id <> 0 AND id IN
	// 						(
	// 							SELECT stale.id
	// 							FROM (
	// 								SELECT id
	// 								FROM Urls
	// 								WHERE last_updated < DATE_SUB( CURRENT_TIMESTAMP(), INTERVAL 6 MONTH)
	// 							) AS stale
	// 						)`

	removeStaleProducts = `DELETE FROM Products WHERE last_updated < NOW() - toIntervalDay(90)`
	removeStaleURLs     = `DELETE FROM Urls WHERE last_updated < NOW() - toIntervalDay(90)`
)

// Cleanup will remove any products and URLs from the database that hasn't been updated
// for more than six months
func Cleanup() error {

	if !viper.GetBool("app.removestale") {
		return errors.Errorf("removal of stale entries is disabled")
	}

	statement, err := GetDB().Prepare(removeStaleProducts)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement for removing stale products")
	}

	res, err := statement.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to remove stale products")
	}

	count, err := res.RowsAffected()
	if err == nil {
		log.Infof("%v stale products removed", count)
	}

	statement, err = GetDB().Prepare(removeStaleURLs)
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement for removing stale URLs")
	}

	res, err = statement.Exec()
	if err != nil {
		return errors.Wrap(err, "failed to remove stale URLs")
	}

	count, err = res.RowsAffected()
	if err == nil {
		log.Infof("%v stale URLs removed", count)
	}

	return nil
}
