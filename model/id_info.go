package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"net/http"
)

type IDInfo struct {
	gorm.Model
	Name   string `gorm:"column:name;size:255;not null;"`
	Number string `gorm:"column:number;size:255;not null;unique;unique_index;"`
}

/*
func FindIDInfoByNumber(number string) (i *IDInfo, err error) {

	i = &IDInfo{}
	if DB.Where("number = ?", number).First(i).RecordNotFound() {
		err = fmt.Errorf("not found")
		return
	}

	if i.DeletedAt != nil {
		err = fmt.Errorf("deleted")
	}

	return
}
*/

func (i *IDInfo) Match() bool {
	return !DB.Where("number = ? AND name = ?", i.Number, i.Name).First(&IDInfo{}).RecordNotFound()
}

func (i *IDInfo) CheckByAPI(conf *config.IDCheckAPI) (ok bool, err error) {

	req := struct {
		Name string `json:"name"`
		ID   string `json:"id_no"`
	}{
		i.Name,
		i.Number,
	}

	b, _ := json.Marshal(req)
	log.Info(string(b))
	log.Info(conf)

	// request instance
	resp, err := resty.R().SetFormData(map[string]string{
		"username": conf.Username,
		"password": fmt.Sprintf("%x", md5.Sum([]byte(conf.Password))),
		"request":  string(b),
	}).
		SetHeader("Accept", "application/json").
		Post(conf.Url)

	if err != nil {
		log.Info(err)
	}
	log.Info(resp)

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("status code: %d", resp.StatusCode())
		return
	}

	result := gjson.GetBytes(resp.Body(), "result")
	if result.String() == "A" {
		ok = true
		DB.Create(i)
	}
	return
}
