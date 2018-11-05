package main

import (
	"github.com/labstack/echo"
	"github.com/lpisces/idcheck/config"
	"github.com/lpisces/idcheck/handler"
	"github.com/lpisces/idcheck/model"
)

func route(e *echo.Echo, conf *config.Config) (err error) {

	if _, err = model.InitDB(conf.DB); err != nil {
		return
	}

	if conf.Debug {
		model.DB.LogMode(true)
	}
	model.Migrate()

	e.GET("/", handler.HandleWelcome())
	e.GET("/id_check", handler.HandleIDCheck(conf.IDCheckAPI))
	e.POST("/upload", handler.HandleIDImageUpload(conf))
	e.POST("/download", handler.HandleIDImageDownload(conf))
	e.POST("/download2", handler.HandleIDImageDownload2(conf))
	e.POST("/sms", handler.HandleSMS(conf.SMSAPI))

	return
}
