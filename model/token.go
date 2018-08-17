package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Name   string `gorm:"column:name;size:255;not null;unique;"`
	Key    string `gorm:"column:key;size:255;not null;unique;unique_index;"`
	Secret string `gorm:"column:secret;size:255;not null;unique;"`
}

func FindTokenByKey(key string) (t *Token, err error) {

	t = &Token{}
	if DB.Unscoped().Where("`key` = ?", key).First(t).RecordNotFound() {
		err = fmt.Errorf("not found")
		return
	}

	if t.DeletedAt != nil {
		err = fmt.Errorf("deleted")
	}

	return
}
