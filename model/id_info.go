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
	Name   string `gorm:"column:name;size:255;not null;unique_index:name_number;index;"`
	Number string `gorm:"column:number;size:255;not null;unique_index:name_number;index;"`
	Result string `gorm:"column:result;size:255;not null;"`
}

func (i *IDInfo) Match() (match bool, cache bool) {
	ii := &IDInfo{}
	if DB.Where("`number` = ? and `name` = ?", i.Number, i.Name).First(ii).RecordNotFound() {
		return false, false
	}
	return ii.Result == "A", true
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

	result := gjson.GetBytes(resp.Body(), "result")

	apiLog := &APILog{
		Request:    string(b),
		StatusCode: resp.StatusCode(),
		Response:   string(resp.Body()),
		Result:     result.String(),
	}
	DB.Create(apiLog)

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("status code: %d", resp.StatusCode())
		return
	}

	i.Result = result.String()
	DB.Create(i)
	if result.String() == "A" {
		ok = true
	}

	return
}
