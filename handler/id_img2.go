package handler

import (
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	//"github.com/lpisces/idcheck/model"
	//"archive/zip"
	"fmt"
	//"github.com/disintegration/imaging"
	"github.com/labstack/gommon/log"
	//"image"
	//"image/color"
	//"io"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleIDImageDownload2(conf *config.Config) func(c echo.Context) error {
	type (
		Ret struct {
			Status   int64
			Msg      string
			NotFound []string
			Url      string
		}

		IDFile struct {
			ID       string
			Filename string
		}
	)

	return func(c echo.Context) error {

		r := Ret{
			0,
			"",
			nil,
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

		ids := c.FormValue("ids")
		path := "./public"

		//id_arr := strings.Split(string(ids), ",")
		var id_map []IDFile
		if err := json.Unmarshal([]byte(ids), &id_map); err != nil {
			log.Info(err)
			return err
		}
		log.Info(id_map)

		var notFound, zipFiles []string
		tmpDir := fmt.Sprintf("/tmp/ids_%s/", time.Now().Format("20060102150405"))
		log.Info(tmpDir)
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			if err := os.Mkdir(tmpDir, os.ModeDir); err != nil {
				return err
			}
		}
		defer os.RemoveAll(tmpDir)
		for _, id := range id_map {
			//log.Info(id)
			filename := conf.IDImageUploadDir + "/" + id.ID + ".jpg"
			newFilename := tmpDir + id.Filename + ".jpg"

			log.Info(newFilename)
			if _, err := os.Stat(newFilename); os.IsNotExist(err) {
				notFound = append(notFound, id.ID)
				continue
			}
			if len(id.Filename) == 0 {
				newFilename = tmpDir + id.ID + ".jpg"
			}

			input, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(newFilename, input, 0644); err != nil {
				return err
			}

			zipFiles = append(zipFiles, newFilename)
		}

		if len(zipFiles) == 0 {
			r.Msg = "no id found"
			r.Status = http.StatusNotFound
			return c.JSON(http.StatusOK, r)
		}

		output := fmt.Sprintf("%s/%s%s.zip", path, "ids_", time.Now().Format("20060102150405"))

		err = ZipFiles(output, zipFiles)
		if err != nil {
			r.Status = http.StatusInternalServerError
			r.Msg = "zip error"
			return c.JSON(http.StatusInternalServerError, r)
		}

		r.Status = http.StatusOK
		r.Msg = "OK"
		r.Url = strings.Replace(output, ".", "", 1)
		r.NotFound = notFound
		return c.JSON(http.StatusOK, r)
	}
}
