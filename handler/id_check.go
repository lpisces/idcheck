package handler

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	"github.com/lpisces/idcheck/model"
	"net/http"
	"strconv"
	"time"
)

var IDCheckHistory map[string]int64
var CheckGapTime int64 = 1000 * 1000 * 1000 // nano second

func init() {
	IDCheckHistory = make(map[string]int64)
}

func HandleIDCheck(conf *config.IDCheckAPI) func(c echo.Context) error {

	type (
		Data struct {
			Cache bool
			Match bool
		}

		Ret struct {
			Status int64
			Msg    string
			*Data
		}
	)

	return func(c echo.Context) error {

		r := Ret{
			0,
			"",
			&Data{
				false,
				false,
			},
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

		i := &model.IDInfo{
			Name:   v.Get("name"),
			Number: v.Get("number"),
		}

		// find in db
		match, cache := i.Match()
		if match {
			r.Status = http.StatusOK
			r.Msg = "OK"
			r.Data.Cache = true
			r.Data.Match = true
			return c.JSON(http.StatusOK, r)
		}
		if !match && cache {
			r.Status = http.StatusOK
			r.Msg = "OK"
			r.Data.Cache = true
			r.Data.Match = false
			return c.JSON(http.StatusOK, r)
		}

		if t, ok := IDCheckHistory[i.Number]; ok {
			if time.Now().UnixNano()-t < CheckGapTime {
				log.Info(i.Number + " waiting")
				time.Sleep(time.Duration(CheckGapTime) * time.Nanosecond)
			}
		}
		IDCheckHistory[i.Number] = time.Now().UnixNano()
		log.Info(fmt.Sprintf("history length is %d", len(IDCheckHistory)))

		// find in api
		ok, err := i.CheckByAPI(conf)
		if err != nil {
			r.Status = http.StatusInternalServerError
			r.Msg = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		if ok {
			r.Status = http.StatusOK
			r.Msg = "OK"
			r.Data.Cache = false
			r.Data.Match = true
			return c.JSON(http.StatusOK, r)
		}

		r.Status = http.StatusOK
		r.Msg = "OK"
		r.Data.Cache = false
		r.Data.Match = false
		return c.JSON(http.StatusOK, r)
	}
}
