package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	"gopkg.in/resty.v1"
	"time"
)

type SMS struct {
	gorm.Model
	Content    string     `gorm:"column:content;size:1024;not null;index;"`
	Mobile     string     `gorm:"column:mobile;size:1024;not null;"`
	TokenID    uint       `gorm:"column:token_id;size:255;not null;"`
	Sign       string     `gorm:"column:sign;size:255;not null;"`
	SendTime   *time.Time `gorm:"column:send_time";not null;`
	Response   string     `gorm:"column:response;not null;"`
	StatusCode string     `gorm:"column:status_code;not null;"`
}

func (sms *SMS) Send(conf *config.SMSAPI) (err error) {

	data := map[string]string{
		"account":  conf.Username,
		"password": conf.Password,
		"mobile":   sms.Mobile,
		"pid":      conf.PID,
		"time":     sms.SendTime.Format("2006-01-02 15:04:05"),
		"content":  fmt.Sprintf("%s【%s】", sms.Content, sms.Sign),
	}

	resp, err := resty.R().SetFormData(data).Post(conf.URL)
	log.Info(resp)

	sms.StatusCode = fmt.Sprintf("%d", resp.StatusCode())
	sms.Response = resp.String()
	DB.Create(sms)

	if err != nil {
		log.Info(err)
		return err
	}

	return
}
