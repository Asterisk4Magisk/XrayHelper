package common

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"encoding/base64"
	"io"
	"os"
	"path"
	"strings"
)

const tagUtil = "util"

// DecodeBase64 decode base64 data
func DecodeBase64(data string) (string, error) {
	decode, err := base64.RawStdEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use RawStdEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.StdEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use StdEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.RawURLEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use RawURLEncoding decode base64 failed, " + err.Error())
	}
	decode, err = base64.URLEncoding.DecodeString(data)
	if err == nil {
		return string(decode), nil
	} else {
		log.HandleDebug("use URLEncoding decode base64 failed, " + err.Error())
	}
	return "", e.New("decode base64 data failed").WithPrefix(tagUtil)
}

// CopyFile copy file from srcName to dstName
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, e.New("open source file failed, ", err).WithPrefix(tagUtil)
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, e.New("create target file failed, ", err).WithPrefix(tagUtil)
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)
	return io.Copy(dst, src)
}

// WildcardMatch simple wildcard matching, time complexity is O(mn)
func WildcardMatch(str string, ptr string) bool {
	if strings.IndexRune(ptr, '*') == -1 && strings.IndexRune(ptr, '?') == -1 {
		return str == ptr
	}
	m, n := len(str), len(ptr)
	dp := make([][]bool, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for i := 1; i <= n; i++ {
		if ptr[i-1] == '*' {
			dp[0][i] = true
		} else {
			break
		}
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if ptr[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if ptr[j-1] == '?' || str[i-1] == ptr[j-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}

// HandleCoreConfDir handle conf dir type json core config like xray, sing-box
func HandleCoreConfDir(handler func(c []byte) (bool, []byte, error)) error {
	confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig)
	if err != nil {
		return e.New("open core config file failed, " + err.Error()).WithPrefix(tagUtil)
	}
	if confInfo.IsDir() {
		if confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig); err == nil {
			for _, conf := range confDir {
				if !conf.IsDir() && strings.HasSuffix(conf.Name(), ".json") {
					if confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name())); err == nil {
						needSave, confByte, err := handler(confByte)
						if err != nil {
							log.HandleDebug(err)
							continue
						}
						if needSave {
							if err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), confByte, 0644); err != nil {
								log.HandleDebug("write new config failed, " + err.Error())
							}
						}
						break
					}
				}
			}
		}
	} else {
		if confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig); err == nil {
			needSave, confByte, err := handler(confByte)
			if err != nil {
				return err
			}
			if needSave {
				if err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig), confByte, 0644); err != nil {
					return e.New("write new config failed, " + err.Error()).WithPrefix(tagUtil)
				}
			}
		} else {
			return e.New("read core config file failed").WithPrefix(tagUtil)
		}
	}
	return nil
}
