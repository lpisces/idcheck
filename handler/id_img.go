package handler

import (
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	//"github.com/lpisces/idcheck/model"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func HandleIDImageUpload(conf *config.Config) func(c echo.Context) error {
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

		id := v.Get("id")
		if len(id) != 18 {
			r.Status = http.StatusUnprocessableEntity
			r.Msg = "id number error"
			return c.JSON(http.StatusUnprocessableEntity, r)
		}

		// Source
		handleUpload := func(filename, storePath string) error {
			file, err := c.FormFile(filename)
			if err != nil {
				return err
			}
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			// Destination
			if _, err := os.Stat(storePath); os.IsNotExist(err) {
				if err := os.Mkdir(storePath, os.ModeDir); err != nil {
					return err
				}
			}
			dst, err := os.Create(storePath + "/" + id + "_" + filename + ".jpg")
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
			return nil
		}

		if err := handleUpload("front", conf.IDImageUploadDir); err != nil {
			r.Status = http.StatusInternalServerError
			r.Msg = "front image upload failed"
			return c.JSON(http.StatusInternalServerError, r)
		}

		if err := handleUpload("back", conf.IDImageUploadDir); err != nil {
			r.Status = http.StatusInternalServerError
			r.Msg = "back image upload failed"
			return c.JSON(http.StatusInternalServerError, r)
		}

		defer func() {
			path1 := conf.IDImageUploadDir + "/" + id + "_" + "front" + ".jpg"
			path2 := conf.IDImageUploadDir + "/" + id + "_" + "back" + ".jpg"
			path3 := conf.IDImageUploadDir + "/" + id + "_" + "merged" + ".jpg"
			mergeImg(path1, path2, path3)
		}()

		r.Status = http.StatusOK
		r.Msg = "OK"
		return c.JSON(http.StatusOK, r)
	}
}

func mergeImg(path1, path2, path3 string) error {
	src1, err := imaging.Open(path1)
	if err != nil {
		return err
	}

	src2, err := imaging.Open(path2)
	if err != nil {
		return err
	}

	src1Fit := imaging.Fit(src1, 500, 300, imaging.Lanczos)
	src2Fit := imaging.Fit(src2, 500, 300, imaging.Lanczos)

	dst := imaging.New(520, 630, color.White)
	dst = imaging.Paste(dst, src1Fit, image.Pt(10, 10))
	dst = imaging.Paste(dst, src2Fit, image.Pt(10, 320))

	return imaging.Save(dst, path3)
}
