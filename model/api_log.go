package model

import (
	//"fmt"
	"github.com/jinzhu/gorm"
)

type APILog struct {
	gorm.Model
	Request    string `gorm:"column:request;size:1024;not null;"`
	StatusCode int    `gorm:"column:status_code;not null;"`
	Response   string `gorm:"column:response;size:1024;not null;"`
	Result     string `gorm:"column:result;size:255;not null;"`
}
