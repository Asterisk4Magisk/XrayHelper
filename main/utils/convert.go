package utils

import (
	"XrayHelper/main/errors"
	"encoding/base64"
)

// DecodeBase64 decode base64 data
func DecodeBase64(data []byte) ([]byte, error) {
	decode, err := base64.URLEncoding.DecodeString(string(data))
	if err != nil {
		return nil, errors.New("decode base64 data failed, ", err).WithPrefix("convert")
	}
	return decode, nil
}
