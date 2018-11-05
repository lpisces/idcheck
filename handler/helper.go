package handler

import (
	"crypto/md5"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/lpisces/idcheck/model"
	"net/url"
)

func checkSign(v url.Values) bool {

	sign := v.Get("sign")
	v.Del("sign")

	key := v.Get("key")
	if "" == key {
		return false
	}
	token, err := model.FindTokenByKey(key)
	if err != nil {
		return false
	}

	orig := token.Secret + v.Encode()
	log.Info(orig)
	vSign := fmt.Sprintf("%x", md5.Sum([]byte(orig)))

	return vSign == sign
}
