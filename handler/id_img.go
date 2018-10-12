package handler

import (
	"github.com/labstack/echo"
	//"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/config"
	//"github.com/lpisces/idcheck/model"
	"archive/zip"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/labstack/gommon/log"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleIDImageDownload(conf *config.Config) func(c echo.Context) error {
	type (
		Ret struct {
			Status   int64
			Msg      string
			NotFound []string
			Url      string
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
		//path := conf.IDImageUploadDir
		path := "./public"

		id_arr := strings.Split(string(ids), ",")

		var notFound, zipFiles []string
		for _, id := range id_arr {
			log.Info(id)
			filename := conf.IDImageUploadDir + "/" + id + ".jpg"
			log.Info(filename)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				notFound = append(notFound, id)
				continue
			}
			zipFiles = append(zipFiles, filename)
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
			path3 := conf.IDImageUploadDir + "/" + id + ".jpg"
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

	watermark, err := imaging.Open("./watermark.png")
	if err != nil {
		return err
	}

	src1Fit := imaging.Fit(src1, 500, 300, imaging.Lanczos)
	src2Fit := imaging.Fit(src2, 500, 300, imaging.Lanczos)
	watermarkFit := imaging.Fit(watermark, 520, 630, imaging.Lanczos)

	dst := imaging.New(520, 630, color.White)
	//dst = imaging.Paste(dst, watermarkFit, image.Pt(0, 0))
	dst = imaging.Paste(dst, src1Fit, image.Pt(10, 10))
	dst = imaging.Paste(dst, src2Fit, image.Pt(10, 320))
	dst = imaging.Overlay(dst, watermarkFit, image.Pt(0, 0), 0.5)

	return imaging.Save(dst, path3)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = file

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
