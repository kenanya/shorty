package helper

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)


func KeyNotFound(err error) bool {
	return strings.Index(err.Error(), "Key not found") >= 0
}

func GetSubStr(value string, from, to int ) string{
	runes := []rune(value)
	return string(runes[from:to])
}

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func DecodeParam(encodedValue string) (string, error) {
	decodedValue, err := url.QueryUnescape(encodedValue)
	if err != nil {
		// log.WithField("error", err).Error("Exception caught")
		// log.Printf("\n DecodeParam : <%+v>\n", err)
		log.Errorf("DecodeParam : %v", err)		
		return "", err
	}
	return decodedValue, nil
}