package handler

import (
	//"fmt"
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	"github.com/lpisces/idcheck/model"
	"net/http"
	"strconv"
	"time"
)

func HandleSMS(conf *config.SMSAPI) func(c echo.Context) error {

	type (
		Ret struct {
			Status int64
			Msg    string
		}
	)

	return func(c echo.Context) error {

		r := Ret{
			0,
			"",
		}

		v := c.QueryParams()
		if !checkSign(v) {
			r.Status = http.StatusUnauthorized
			r.Msg = "invalid sign"
			return c.JSON(http.StatusUnauthorized, r)
		}

		expire, err := strconv.ParseInt(v.Get("expire"), 10, 64)
		if err != nil {
			return err
		}
		if expire < time.Now().Unix() {
			r.Status = http.StatusUnauthorized
			r.Msg = "expired"
			return c.JSON(http.StatusUnauthorized, r)
		}

		sms := &model.SMS{
			Sign:    v.Get("sms_sign"),
			Content: v.Get("sms_content"),
			Mobile:  v.Get("sms_mobile"),
		}

		st := time.Now().Add(time.Second * 5)
		sms.SendTime = &st

		key := v.Get("key")
		token, _ := model.FindTokenByKey(key)

		sms.TokenID = token.ID
		if err := sms.Send(conf); err != nil {
			r.Status = http.StatusInternalServerError
			r.Msg = err.Error()
			return c.JSON(http.StatusInternalServerError, r)
		}

		r.Status = http.StatusOK
		r.Msg = "OK"
		return c.JSON(http.StatusOK, r)
	}
}
