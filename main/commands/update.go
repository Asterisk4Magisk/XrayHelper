package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	singboxUrl                = "https://api.github.com/repos/SagerNet/sing-box/releases/latest"
	clashUrl                  = "https://api.github.com/repos/Dreamacro/clash/releases/latest"
	yacdDownloadUrl           = "https://github.com/haishanh/yacd/archive/gh-pages.zip"
	xrayCoreDownloadUrl       = "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-android-arm64-v8a.zip"
	v2rayCoreDownloadUrl      = "https://github.com/v2fly/v2ray-core/releases/latest/download/v2ray-android-arm64-v8a.zip"
	geoipDownloadUrl          = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
	geoipDownloadUrlSingbox   = "https://github.com/lyc8503/sing-box-rules/releases/latest/download/geoip.db"
	geositeDownloadUrl        = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
	geositeDownloadUrlSingbox = "https://github.com/lyc8503/sing-box-rules/releases/latest/download/geosite.db"
	tun2socksDownloadUrl      = "https://github.com/heiher/hev-socks5-tunnel/releases/latest/download/hev-socks5-tunnel-linux-arm64"
)

type UpdateCommand struct{}

func (this *UpdateCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify operation, available operation [core|tun2socks|geodata|subscribe]").WithPrefix("update").WithPathObj(*this)
	}
	if len(args) > 1 {
		return errors.New("too many arguments").WithPrefix("update").WithPathObj(*this)
	}
	switch args[0] {
	case "core":
		log.HandleInfo("update: updating core")
		if err := updateCore(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	case "tun2socks":
		log.HandleInfo("update: updating tun2socks")
		if err := updateTun2socks(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	case "geodata":
		log.HandleInfo("update: updating geodata")
		if err := updateGeodata(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	case "subscribe":
		log.HandleInfo("update: updating subscribe")
		if err := updateSubscribe(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	case "yacd":
		log.HandleInfo("update: updating yacd")
		if err := updateYacd(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [core|tun2socks|geodata|subscribe|yacd]").WithPrefix("update").WithPathObj(*this)
	}
	return nil
}

// updateCore update core, support xray, singbox
func updateCore() error {
	if runtime.GOARCH != "arm64" {
		return errors.New("this feature only support arm64 device").WithPrefix("update")
	}
	serviceRunFlag := false
	if err := os.MkdirAll(builds.Config.XrayHelper.DataDir, 0644); err != nil {
		return errors.New("create run dir failed, ", err).WithPrefix("update")
	}
	if err := os.MkdirAll(path.Dir(builds.Config.XrayHelper.CorePath), 0644); err != nil {
		return errors.New("create core path dir failed, ", err).WithPrefix("update")
	}
	switch builds.Config.XrayHelper.CoreType {
	case "xray":
		xrayZipPath := path.Join(builds.Config.XrayHelper.DataDir, "xray.zip")
		if err := common.DownloadFile(xrayZipPath, xrayCoreDownloadUrl); err != nil {
			return err
		}
		// update core need stop core service first
		if len(getServicePid()) > 0 {
			log.HandleInfo("update: detect core is running, stop it")
			stopService()
			serviceRunFlag = true
			_ = os.Remove(builds.Config.XrayHelper.CorePath)
		}
		zipReader, err := zip.OpenReader(xrayZipPath)
		if err != nil {
			return errors.New("open xray.zip failed, ", err).WithPrefix("update")
		}
		defer func(zipReader *zip.ReadCloser) {
			_ = zipReader.Close()
			_ = os.Remove(xrayZipPath)
		}(zipReader)
		for _, file := range zipReader.File {
			if file.Name == "xray" {
				fileReader, err := file.Open()
				if err != nil {
					return errors.New("cannot get file reader "+file.Name+", ", err).WithPrefix("update")
				}
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				if err != nil {
					return errors.New("cannot open file "+builds.Config.XrayHelper.CorePath+", ", err).WithPrefix("update")
				}
				_, err = io.Copy(saveFile, fileReader)
				if err != nil {
					return errors.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
				}
				_ = saveFile.Close()
				_ = fileReader.Close()
				break
			}
		}
	case "v2ray":
		v2rayZipPath := path.Join(builds.Config.XrayHelper.DataDir, "v2ray.zip")
		if err := common.DownloadFile(v2rayZipPath, v2rayCoreDownloadUrl); err != nil {
			return err
		}
		// update core need stop core service first
		if len(getServicePid()) > 0 {
			log.HandleInfo("update: detect core is running, stop it")
			stopService()
			serviceRunFlag = true
			_ = os.Remove(builds.Config.XrayHelper.CorePath)
		}
		zipReader, err := zip.OpenReader(v2rayZipPath)
		if err != nil {
			return errors.New("open v2ray.zip failed, ", err).WithPrefix("update")
		}
		defer func(zipReader *zip.ReadCloser) {
			_ = zipReader.Close()
			_ = os.Remove(v2rayZipPath)
		}(zipReader)
		for _, file := range zipReader.File {
			if file.Name == "v2ray" {
				fileReader, err := file.Open()
				if err != nil {
					return errors.New("cannot get file reader "+file.Name+", ", err).WithPrefix("update")
				}
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				if err != nil {
					return errors.New("cannot open file "+builds.Config.XrayHelper.CorePath+", ", err).WithPrefix("update")
				}
				_, err = io.Copy(saveFile, fileReader)
				if err != nil {
					return errors.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
				}
				_ = saveFile.Close()
				_ = fileReader.Close()
				break
			}
		}
	case "sing-box":
		singboxDownloadUrl, err := getDownloadUrl(singboxUrl, "android-arm64.tar.gz")
		if err != nil {
			return err
		}
		singboxGzipPath := path.Join(builds.Config.XrayHelper.DataDir, "sing-box.tar.gz")
		if err := common.DownloadFile(singboxGzipPath, singboxDownloadUrl); err != nil {
			return err
		}
		// update core need stop core service first
		if len(getServicePid()) > 0 {
			log.HandleInfo("update: detect core is running, stop it")
			stopService()
			serviceRunFlag = true
			_ = os.Remove(builds.Config.XrayHelper.CorePath)
		}
		singboxGzip, err := os.Open(singboxGzipPath)
		if err != nil {
			return errors.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(singboxGzip *os.File) {
			_ = singboxGzip.Close()
			_ = os.Remove(singboxGzipPath)
		}(singboxGzip)
		gzipReader, err := gzip.NewReader(singboxGzip)
		if err != nil {
			return errors.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(gzipReader *gzip.Reader) {
			_ = gzipReader.Close()
		}(gzipReader)
		tarReader := tar.NewReader(gzipReader)
		for {
			fileHeader, err := tarReader.Next()
			if err != nil {
				if err == io.EOF {
					return errors.New("cannot find sing-box binary").WithPrefix("update")
				}
				continue
			}
			if filepath.Base(fileHeader.Name) == "sing-box" {
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				_, err = io.Copy(saveFile, tarReader)
				if err != nil {
					return errors.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
				}
				_ = saveFile.Close()
				break
			}
		}
	case "clash":
		clashDownloadUrl, err := getDownloadUrl(clashUrl, "clash-linux-arm64")
		if err != nil {
			return err
		}
		clashGzipPath := path.Join(builds.Config.XrayHelper.DataDir, "clash.gz")
		if err := common.DownloadFile(clashGzipPath, clashDownloadUrl); err != nil {
			return err
		}
		// update core need stop core service first
		if len(getServicePid()) > 0 {
			log.HandleInfo("update: detect core is running, stop it")
			stopService()
			serviceRunFlag = true
			_ = os.Remove(builds.Config.XrayHelper.CorePath)
		}
		clashGzip, err := os.Open(clashGzipPath)
		if err != nil {
			return errors.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(clashGzip *os.File) {
			_ = clashGzip.Close()
			_ = os.Remove(clashGzipPath)
		}(clashGzip)
		gzipReader, err := gzip.NewReader(clashGzip)
		if err != nil {
			return errors.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(gzipReader *gzip.Reader) {
			_ = gzipReader.Close()
		}(gzipReader)
		saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
		_, err = io.Copy(saveFile, gzipReader)
		if err != nil {
			return errors.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
		}
		_ = saveFile.Close()
	default:
		return errors.New("unknown core type " + builds.Config.XrayHelper.CoreType).WithPrefix("update")
	}
	if serviceRunFlag {
		log.HandleInfo("update: starting core with new version")
		_ = startService()
	}
	return nil
}

// updateTun2socks update tun2socks
func updateTun2socks() error {
	if runtime.GOARCH != "arm64" {
		return errors.New("this feature only support arm64 device").WithPrefix("update")
	}
	savePath := path.Join(path.Dir(builds.Config.XrayHelper.CorePath), "tun2socks")
	if err := common.DownloadFile(savePath, tun2socksDownloadUrl); err != nil {
		return err
	}
	return nil
}

// updateGeodata update geodata
func updateGeodata() error {
	if err := os.MkdirAll(builds.Config.XrayHelper.DataDir, 0644); err != nil {
		return errors.New("create DataDir failed, ", err).WithPrefix("update")
	}
	switch builds.Config.XrayHelper.CoreType {
	case "sing-box":
		if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geoip.db"), geoipDownloadUrlSingbox); err != nil {
			return err
		}
		if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geosite.db"), geositeDownloadUrlSingbox); err != nil {
			return err
		}
	default:
		if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geoip.dat"), geoipDownloadUrl); err != nil {
			return err
		}
		if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geosite.dat"), geositeDownloadUrl); err != nil {
			return err
		}
	}
	return nil
}

// updateSubscribe update subscribe
func updateSubscribe() error {
	if err := os.MkdirAll(builds.Config.XrayHelper.DataDir, 0644); err != nil {
		return errors.New("create DataDir failed, ", err).WithPrefix("update")
	}
	if builds.Config.XrayHelper.CoreType == "clash" {
		for index, subUrl := range builds.Config.XrayHelper.SubList {
			rawData, err := common.GetRawData(subUrl)
			if err != nil {
				log.HandleError(err)
				continue
			}
			if err := os.WriteFile(path.Join(builds.Config.XrayHelper.DataDir, "clashSub"+strconv.Itoa(index)+".yaml"), rawData, 0644); err != nil {
				return errors.New("write subscribe file failed, ", err).WithPrefix("update")
			}
		}
	} else {
		builder := strings.Builder{}
		for _, subUrl := range builds.Config.XrayHelper.SubList {
			rawData, err := common.GetRawData(subUrl)
			if err != nil {
				log.HandleError(err)
				continue
			}
			subData, err := common.DecodeBase64(string(rawData))
			if err != nil {
				log.HandleError(err)
				continue
			}
			builder.WriteString(strings.TrimSpace(subData) + "\n")
		}
		if builder.Len() > 0 {
			if err := os.WriteFile(path.Join(builds.Config.XrayHelper.DataDir, "sub.txt"), []byte(builder.String()), 0644); err != nil {
				return errors.New("write subscribe file failed, ", err).WithPrefix("update")
			}
		}
	}
	return nil
}

// updateYacd update yacd
func updateYacd() error {
	yacdZipPath := path.Join(builds.Config.XrayHelper.DataDir, "yacd.zip")
	if err := common.DownloadFile(yacdZipPath, yacdDownloadUrl); err != nil {
		return err
	}
	zipReader, err := zip.OpenReader(yacdZipPath)
	if err != nil {
		return errors.New("open yacd.zip failed, ", err).WithPrefix("update")
	}
	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
		_ = os.Remove(yacdZipPath)
	}(zipReader)
	if err := os.RemoveAll(builds.Config.XrayHelper.DataDir); err != nil {
		return errors.New("remove old yacd files failed, ", err).WithPrefix("update")
	}
	for _, file := range zipReader.File {
		t := filepath.Join(builds.Config.XrayHelper.DataDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(t, 0644); err != nil {
				return errors.New("create dir "+t+" failed, ", err).WithPrefix("update")
			}
			continue
		}
		fr, err := file.Open()
		if err != nil {
			return errors.New("open file "+file.Name+" failed, ", err)
		}
		fw, err := os.OpenFile(t, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			_ = fr.Close()
			return errors.New("open file "+t+" failed, ", err).WithPrefix("update")
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			_ = fw.Close()
			_ = fr.Close()
			return errors.New("copy file "+file.Name+" failed, ", err).WithPrefix("update")
		}
		_ = fw.Close()
		_ = fr.Close()
	}
	return nil
}

// getDownloadUrl use github api to get download url
func getDownloadUrl(githubApi string, nameContent string) (string, error) {
	rawData, err := common.GetRawData(githubApi)
	if err != nil {
		return "", err
	}
	var jsonValue interface{}
	err = json.Unmarshal(rawData, &jsonValue)
	if err != nil {
		return "", errors.New("unmarshal github json failed, ", err).WithPrefix("update")
	}
	// assert json to map
	jsonMap, ok := jsonValue.(map[string]interface{})
	if !ok {
		return "", errors.New("assert github json to map failed").WithPrefix("update")
	}
	assets, ok := jsonMap["assets"]
	if !ok {
		return "", errors.New("cannot find assets ").WithPrefix("update")
	}
	// assert assets
	assetsMap, ok := assets.([]interface{})
	if !ok {
		return "", errors.New("assert assets to []interface failed").WithPrefix("update")
	}
	for _, asset := range assetsMap {
		assetMap, ok := asset.(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := assetMap["name"].(string)
		if !ok {
			continue
		}
		if strings.Contains(name, nameContent) {
			downloadUrl, ok := assetMap["browser_download_url"].(string)
			if !ok {
				return "", errors.New("assert browser_download_url to string failed").WithPrefix("update")
			}
			return downloadUrl, nil
		}
	}
	return "", errors.New("cannot get download url from " + githubApi).WithPrefix("update")
}
