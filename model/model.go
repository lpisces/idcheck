package model

import (
	//"github.com/labstack/gommon/log"

	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lpisces/idcheck/config"
	"time"
)

func InitDB(conf *config.DB) (*gorm.DB, error) {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Database))
	if err != nil {
		return nil, err
	}
	DB = db
	Config = conf

	go func(db *gorm.DB) {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			<-ticker.C
			db.DB().Ping()
		}
	}(db)
	return db, err
}

func GetDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	if nil == Config {
		err := fmt.Errorf("no db config")
		return nil, err
	}

	// reconnect
	return InitDB(Config)
}

func Migrate() {
	DB.AutoMigrate(&Token{}, &IDInfo{}, &APILog{}, &SMS{})
}
