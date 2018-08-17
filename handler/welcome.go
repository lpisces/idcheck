package handler

import (
	//"fmt"
	"github.com/labstack/echo"
	//"github.com/lpisces/idcheck/config"
	"net/http"
)

func HandleWelcome() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "welcome")
	}
}
