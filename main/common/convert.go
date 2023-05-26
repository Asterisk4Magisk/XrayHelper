package common

import (
	"XrayHelper/main/errors"
	"encoding/base64"
)

// DecodeBase64 decode base64 data
func DecodeBase64(data string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("decode base64 data failed, ", err).WithPrefix("convert")
	}
	return string(decode), nil
}
