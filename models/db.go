package models

import (
	"database/sql"
	stdlog "log"
	"os"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"

	"github.com/apex/log"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db     *sql.DB
	gormDB *gorm.DB
)

func Connect() *sql.DB {
	var err error
	dsn := "clickhouse://130.211.211.73:9000/?user=default&password=default&database=retroced"
	// dsn := "clickhouse://localhost:9000/?user=&password=&database=retroced"

	Db, err = sql.Open("clickhouse", dsn)
	//Db, err = sql.Open("mysql", fmt.Sprintf(
	//	"%s:%s@tcp(%s:%d)/%s?parseTime=true",
	//	viper.GetString("mysql.user"),
	//	viper.GetString("mysql.pass"),
	//	viper.GetString("mysql.host"),
	//	viper.GetInt("mysql.port"),
	//	viper.GetString("mysql.db"),
	//))
	// To make sure that we don't get dropped requests
	// on slow days.
	if err != nil {
		log.WithError(err).Fatal("error connecting to DB")
	}
	Db.SetMaxIdleConns(0)

	err = Db.Ping()
	if err != nil {
		log.WithError(err).Fatal("error while pinging DB")
	}

	gormLogger := logger.New(
		stdlog.New(os.Stderr, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)
	gormDB, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.WithError(err).Fatal("error creating gorm DB connection")
	}
	gormDB, err = gorm.Open(mysql.New(mysql.Config{Conn: Db}), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.WithError(err).Fatal("error creating gorm DB connection")
	}
	//gormDB.AutoMigrate(&Site{})
	//gormDB.AutoMigrate(&Click{})
	//gormDB.AutoMigrate(&Favourites{})
	//gormDB.AutoMigrate(&ProductRecord{})
	//gormDB.AutoMigrate(&Url{})
	//gormDB.AutoMigrate(&User{})

	return Db
}

// GetDB returns a handle to the database
func GetDB() *sql.DB {
	return Db
}

// GetDBv2 returns a Gorm database connection handle. In the future, this should replace the GetDB() functionality.
func GetDBv2() *gorm.DB {
	return gormDB
}
