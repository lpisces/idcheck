package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lpisces/idcheck/config"
)

var (
	DB     *gorm.DB
	Config *config.DB
)
