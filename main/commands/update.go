package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxies"
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
	mihomoUrl                 = "https://api.github.com/repos/MetaCubeX/mihomo/releases/latest"
	yacdMetaDownloadUrl       = "https://github.com/MetaCubeX/yacd/archive/gh-pages.zip"
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
		return e.New("not specify operation, available operation [core|tun2socks|geodata|subscribe|yacd-meta]").WithPrefix("update").WithPathObj(*this)
	}
	if len(args) > 1 {
		return e.New("too many arguments").WithPrefix("update").WithPathObj(*this)
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
	case "yacd-meta":
		log.HandleInfo("update: updating yacd-meta")
		if err := updateYacdMeta(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	default:
		return e.New("unknown operation " + args[0] + ", available operation [core|tun2socks|geodata|subscribe|yacd|yacd-meta]").WithPrefix("update").WithPathObj(*this)
	}
	return nil
}

// updateCore update core, support xray, singbox
func updateCore() error {
	if runtime.GOARCH != "arm64" {
		return e.New("this feature only support arm64 device").WithPrefix("update")
	}
	serviceRunFlag := false
	if err := os.MkdirAll(builds.Config.XrayHelper.DataDir, 0644); err != nil {
		return e.New("create run dir failed, ", err).WithPrefix("update")
	}
	if err := os.MkdirAll(path.Dir(builds.Config.XrayHelper.CorePath), 0644); err != nil {
		return e.New("create core path dir failed, ", err).WithPrefix("update")
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
			return e.New("open xray.zip failed, ", err).WithPrefix("update")
		}
		defer func(zipReader *zip.ReadCloser) {
			_ = zipReader.Close()
			_ = os.Remove(xrayZipPath)
		}(zipReader)
		for _, file := range zipReader.File {
			if file.Name == "xray" {
				fileReader, err := file.Open()
				if err != nil {
					return e.New("cannot get file reader "+file.Name+", ", err).WithPrefix("update")
				}
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				if err != nil {
					return e.New("cannot open file "+builds.Config.XrayHelper.CorePath+", ", err).WithPrefix("update")
				}
				_, err = io.Copy(saveFile, fileReader)
				if err != nil {
					return e.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
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
			return e.New("open v2ray.zip failed, ", err).WithPrefix("update")
		}
		defer func(zipReader *zip.ReadCloser) {
			_ = zipReader.Close()
			_ = os.Remove(v2rayZipPath)
		}(zipReader)
		for _, file := range zipReader.File {
			if file.Name == "v2ray" {
				fileReader, err := file.Open()
				if err != nil {
					return e.New("cannot get file reader "+file.Name+", ", err).WithPrefix("update")
				}
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				if err != nil {
					return e.New("cannot open file "+builds.Config.XrayHelper.CorePath+", ", err).WithPrefix("update")
				}
				_, err = io.Copy(saveFile, fileReader)
				if err != nil {
					return e.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
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
			return e.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(singboxGzip *os.File) {
			_ = singboxGzip.Close()
			_ = os.Remove(singboxGzipPath)
		}(singboxGzip)
		gzipReader, err := gzip.NewReader(singboxGzip)
		if err != nil {
			return e.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(gzipReader *gzip.Reader) {
			_ = gzipReader.Close()
		}(gzipReader)
		tarReader := tar.NewReader(gzipReader)
		for {
			fileHeader, err := tarReader.Next()
			if err != nil {
				if err == io.EOF {
					return e.New("cannot find sing-box binary").WithPrefix("update")
				}
				continue
			}
			if filepath.Base(fileHeader.Name) == "sing-box" {
				saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
				_, err = io.Copy(saveFile, tarReader)
				if err != nil {
					return e.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
				}
				_ = saveFile.Close()
				break
			}
		}
	case "clash.meta", "mihomo":
		mihomoDownloadUrl, err := getDownloadUrl(mihomoUrl, "mihomo-android-arm64-v")
		if err != nil {
			return err
		}
		mihomoGzipPath := path.Join(builds.Config.XrayHelper.DataDir, "mihomo.gz")
		if err := common.DownloadFile(mihomoGzipPath, mihomoDownloadUrl); err != nil {
			return err
		}
		// update core need stop core service first
		if len(getServicePid()) > 0 {
			log.HandleInfo("update: detect core is running, stop it")
			stopService()
			serviceRunFlag = true
			_ = os.Remove(builds.Config.XrayHelper.CorePath)
		}
		mihomoGzip, err := os.Open(mihomoGzipPath)
		if err != nil {
			return e.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(mihomoGzip *os.File) {
			_ = mihomoGzip.Close()
			_ = os.Remove(mihomoGzipPath)
		}(mihomoGzip)
		gzipReader, err := gzip.NewReader(mihomoGzip)
		if err != nil {
			return e.New("open gzip file failed, ", err).WithPrefix("update")
		}
		defer func(gzipReader *gzip.Reader) {
			_ = gzipReader.Close()
		}(gzipReader)
		saveFile, err := os.OpenFile(builds.Config.XrayHelper.CorePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
		_, err = io.Copy(saveFile, gzipReader)
		if err != nil {
			return e.New("save file "+builds.Config.XrayHelper.CorePath+" failed, ", err).WithPrefix("update")
		}
		_ = saveFile.Close()
	default:
		return e.New("unknown core type " + builds.Config.XrayHelper.CoreType).WithPrefix("update")
	}
	if serviceRunFlag {
		log.HandleInfo("update: starting core with new version")
		_ = startService()
		if err := builds.LoadPackage(); err != nil {
			log.HandleError("update: load package failed, " + err.Error())
		} else {
			proxy, err := proxies.NewProxy(builds.Config.Proxy.Method)
			if err != nil {
				log.HandleError("update: get proxy failed, " + err.Error())
			} else {
				proxy.Disable()
				_ = proxy.Enable()
			}
		}
	}
	return nil
}

// updateTun2socks update tun2socks
func updateTun2socks() error {
	if runtime.GOARCH != "arm64" {
		return e.New("this feature only support arm64 device").WithPrefix("update")
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
		return e.New("create DataDir failed, ", err).WithPrefix("update")
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
		return e.New("create DataDir failed, ", err).WithPrefix("update")
	}
	var v2rayNgUrl, clashUrl []string
	for _, subUrl := range builds.Config.XrayHelper.SubList {
		if strings.HasPrefix(subUrl, "clash+") {
			clashUrl = append(clashUrl, strings.TrimPrefix(subUrl, "clash+"))
		} else {
			v2rayNgUrl = append(v2rayNgUrl, subUrl)
		}
	}
	// update v2rayNg subscribe
	builder := strings.Builder{}
	for _, subUrl := range v2rayNgUrl {
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
			return e.New("write subscribe file failed, ", err).WithPrefix("update")
		}
	}
	// update clash subscribe
	for index, subUrl := range clashUrl {
		rawData, err := common.GetRawData(subUrl)
		if err != nil {
			log.HandleError(err)
			continue
		}
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.DataDir, "clashSub"+strconv.Itoa(index)+".yaml"), rawData, 0644); err != nil {
			return e.New("write subscribe file failed, ", err).WithPrefix("update")
		}
	}
	return nil
}

// updateYacdMeta update yacd-meta
func updateYacdMeta() error {
	yacdMetaZipPath := path.Join(builds.Config.XrayHelper.DataDir, "yacd-meta.zip")
	if err := common.DownloadFile(yacdMetaZipPath, yacdMetaDownloadUrl); err != nil {
		return err
	}
	zipReader, err := zip.OpenReader(yacdMetaZipPath)
	if err != nil {
		return e.New("open zip file failed, ", err).WithPrefix("update")
	}
	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
		_ = os.Remove(yacdMetaZipPath)
	}(zipReader)
	if err := os.RemoveAll(path.Join(builds.Config.XrayHelper.DataDir, "Yacd-meta-gh-pages/")); err != nil {
		return e.New("remove old yacd-meta files failed, ", err).WithPrefix("update")
	}
	for _, file := range zipReader.File {
		t := filepath.Join(builds.Config.XrayHelper.DataDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(t, 0644); err != nil {
				return e.New("create dir "+t+" failed, ", err).WithPrefix("update")
			}
			continue
		}
		fr, err := file.Open()
		if err != nil {
			return e.New("open file "+file.Name+" failed, ", err)
		}
		fw, err := os.OpenFile(t, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			_ = fr.Close()
			return e.New("open file "+t+" failed, ", err).WithPrefix("update")
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			_ = fw.Close()
			_ = fr.Close()
			return e.New("copy file "+file.Name+" failed, ", err).WithPrefix("update")
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
		return "", e.New("unmarshal github json failed, ", err).WithPrefix("update")
	}
	// assert json to map
	jsonMap, ok := jsonValue.(map[string]interface{})
	if !ok {
		return "", e.New("assert github json to map failed").WithPrefix("update")
	}
	assets, ok := jsonMap["assets"]
	if !ok {
		return "", e.New("cannot find assets ").WithPrefix("update")
	}
	// assert assets
	assetsMap, ok := assets.([]interface{})
	if !ok {
		return "", e.New("assert assets to []interface failed").WithPrefix("update")
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
				return "", e.New("assert browser_download_url to string failed").WithPrefix("update")
			}
			return downloadUrl, nil
		}
	}
	return "", e.New("cannot get download url from " + githubApi).WithPrefix("update")
}
