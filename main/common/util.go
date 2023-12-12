package common

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"encoding/base64"
	"io"
	"os"
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
	return "", e.New("decode base64 data failed").WithPrefix("util")
}

// CopyFile copy file from srcName to dstName
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, e.New("open source file failed, ", err).WithPrefix("util")
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, e.New("create target file failed, ", err).WithPrefix("util")
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)
	return io.Copy(dst, src)
}
