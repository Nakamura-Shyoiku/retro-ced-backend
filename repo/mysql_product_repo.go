package repo

import (
	"fmt"
	"github.com/ulventech/retro-ced-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlProductRepo struct {
	DB *gorm.DB
}

func NewMysqlProductRepo(conn string)(*MysqlProductRepo, error) {
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s err %s", conn, err)
	}
	return &MysqlProductRepo{
		DB: db,
	},nil
}

func (m *MysqlProductRepo) InsertProducts(p []models.Product) error {
	r := m.DB.Table("Products").Create(&p)
	if r.Error  != nil {
		return fmt.Errorf("unable to insert products %v due to %s",p, r.Error)
	}
	return nil
}

