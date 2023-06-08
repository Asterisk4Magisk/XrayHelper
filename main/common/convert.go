package common

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"encoding/base64"
)

// DecodeBase64 decode base64 data
func DecodeBase64(data string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use StdEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.RawStdEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use RawStdEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.URLEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use URLEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.RawURLEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use RawURLEncoding decode base64 failed, " + err.Error())
	}
	return "", errors.New("decode base64 data failed").WithPrefix("convert")
}
