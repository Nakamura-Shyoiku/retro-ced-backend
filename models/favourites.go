package models

import (
	"database/sql"

	"github.com/apex/log"
)

type Favourites struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`
	UserId    int64 `json:"user_id"`
}

func (f *Favourites) TableName() string {
	return "Favourites";
}

func (f *Favourites) AddFavourite() error {
	stmt, err := GetDB().Prepare(`
		INSERT INTO Favourites(
			product_id,
			user_id
		) VALUES (?, ?)
	`)
	if err != nil {
		log.WithError(err).Error("failed prepare insert Favourites statement")
		return err
	}
	_, err = stmt.Exec(
		f.ProductId,
		f.UserId,
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on insert Favourites")
		return err
	}

	return nil
}

func (f *Favourites) RemoveFavourite() error {
	stmt, err := GetDB().Prepare(`
	DELETE FROM Favourites
	WHERE product_id = ?
	AND user_id = ?
	`)
	if err != nil {
		log.WithError(err).Error("failed prepare delete Favourites statement")
		return err
	}
	_, err = stmt.Exec(
		f.ProductId,
		f.UserId,
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on delete Favourites")
		return err
	}

	return nil
}

func (f *Favourites) GetFavourite(pid, uid int64) (Favourites, error) {
	var fav Favourites
	err := GetDB().QueryRow(`
		SELECT id, product_id, user_id
		FROM Favourites
		WHERE product_id = ?
		AND user_id = ?`,
		pid,
		uid,
	).Scan(
		&fav.Id,
		&fav.ProductId,
		&fav.UserId,
	)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("fetch Favourite failed")
		return fav, err
	}

	return fav, nil
}
