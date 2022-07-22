package models

import (
	"database/sql"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
)

// User model
type User struct {
	Id            int64     `json:"id"`
	Guid          string    `json:"guid"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	FirstTime     bool      `json:"first_time"`
	FbID          string    `json:"fb_id"`
	CreatedAt     time.Time `json:"created_at"`
	ACL           int       `json:"acl"`
	PartnerSiteID int64     `json:"partner_site_id"`
	PasswordReset string    `json:"password_reset"`
	// Password should NEVER be returned in a json response
	Password string `json:"-"`
}

func (u *User) TableName() string {
	return "Users"
}

// GetUserByGUID returns user by guid
func GetUserByGUID(guid string) (User, error) {
	var u User
	err := GetDBv2().Table("Users").Find(&u, "guid = ?", guid).Error
	if err != nil {
		log.WithError(err).Error("fetch user failed")
		return User{}, err
	}
	return u, nil
}

// GetUserByEmail returns user by email
func GetUserByEmail(email string) (User, error) {
	var u User
	err := GetDBv2().Table("Users").Find(&u, "email = ?", email).Error
	if err != nil {
		log.WithError(err).Error("fetch user failed")
		return User{}, err
	}
	return u, nil
}

// GetUserByID returns user by Guid
func (u *User) GetUserByID(userID int64) (*User, error) {
	var res User
	err := GetDBv2().Table("Users").First(&res, "id = ?", userID).Error
	if err != nil {
		log.WithError(err).Error("fetch user failed")
		return nil, err
	}
	return &res, nil

}

// Create inserts new user
func (u *User) Create() (guid string, err error) {
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		log.WithError(err).Error("failed begin insert transaction ")
		return "", err
	}

	var autoIncresedId int64

	err = GetDB().QueryRow(`
		SELECT COALESCE(MAX(id) + 1, 0)
		FROM Users
		`).Scan(
		&autoIncresedId,
	)
	if err != nil {
		log.WithError(err).Error("failed get max user Id ")
		return "", err
	}

	stmt, err := tx.Prepare("INSERT INTO Users(id, guid, email, password, first_name, last_name, fb_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.WithError(err).Error("failed prepare insert user statement")
		return "", err
	}

	guid = uuid.New().String()
	_, err = stmt.Exec(
		autoIncresedId,
		guid,
		u.Email,
		u.Password,
		u.FirstName,
		u.LastName,
		u.FbID,
		time.Now(),
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on insert user")
		return "", err
	}

	if err := tx.Commit(); err != nil {
		log.WithError(err).Error("failed to commit insert")
		return "", err
	}

	return guid, nil
}

// UpdatePasswordReset updates the password reset hash
func (u *User) UpdatePasswordReset(token string, email string) error {

	stmt, err := GetDB().Prepare("ALTER TABLE Users UPDATE password_reset = ? WHERE email = ?")
	if err != nil {
		log.WithError(err).Error("failed prepare update user statement")
		return err
	}

	_, err = stmt.Exec(token, email)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update user")
		return err
	}

	return nil
}

// UpdateUserPass updates the users password
func (u *User) UpdateUserPass(pass, email string) error {
	stmt, err := GetDB().Prepare("ALTER TABLE Users UPDATE password = ? WHERE email = ?")
	if err != nil {
		log.WithError(err).Error("failed prepare update user statement")
		return err
	}
	_, err = stmt.Exec(pass, email)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update user")
		return err
	}

	return nil
}

// UpdateFbID updates users facebook Guid
func UpdateFbID(guid, fbID string) error {
	stmt, err := GetDB().Prepare("ALTER TABLE Users UPDATE Users SET fb_id = ? WHERE guid = ?")
	if err != nil {
		log.WithError(err).Error("failed prepare update user statement")
		return err
	}
	_, err = stmt.Exec(fbID, guid)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update user")
		return err
	}
	return nil
}

// UpdateUserACL updates permission and siteID (for PARTNER type) of a user.
func (u *User) UpdateUserACL(userID int64, acl int, siteID int64) error {
	stmt, err := GetDB().Prepare(
		"ALTER TABLE Users  UPDATE acl = ?, partner_site_id = ? WHERE id = ?",
	)
	if err != nil {
		log.WithError(err).Error("failed to prepare update user statement")
		return err
	}
	_, err = stmt.Exec(
		acl,
		siteID,
		userID,
	)
	if err != nil {
		log.WithError(err).Error("failed to update user")
		return err
	}

	return nil
}

// GetUsers gets all users if searchquery is null.
func GetUsers(email string, offset int, limit int) ([]User, error) {
	var users []User

	db := GetDBv2().
		Table("Users").
		Offset(offset).
		Limit(limit).
		Order("created_at desc")

	if email != "" {
		term := "%" + email + "%"
		db = db.Where("email like ?", term)
	}
	err := db.Scan(&users).Error

	if err != nil {
		return users, err
	}

	return users, nil
}

// GetUsersCount get the number of users for pagination purposes.
func GetUsersCount(search string) (count uint64, err error) {
	searchQuery := "%" + search + "%"
	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM Users WHERE email LIKE ? or first_name LIKE ? or last_name LIKE ?`,
		searchQuery,
		searchQuery,
		searchQuery,
	).Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("users count row error")
		return count, err
	}

	return count, nil
}

// VerifyPermission verifies whether the user has the required permission
func VerifyPermission(userID string, permissionID int64) error {
	var user User
	var err error

	err = GetDB().QueryRow(`
		SELECT id FROM Users
		WHERE guid = ?
		AND (acl = ? OR acl = 1000)
		`,
		userID,
		permissionID,
	).Scan(
		&user.Id,
	)
	if err != nil && err == sql.ErrNoRows {
		log.WithError(err).Error("No User record")
		return err
	}
	if err != nil {
		log.WithError(err).Error("Verifying user permission failed")
		return err
	}

	return nil
}
