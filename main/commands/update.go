package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"archive/zip"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	xrayCoreDownloadUrl  = "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-android-arm64-v8a.zip"
	geoipDownloadUrl     = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
	geositeDownloadUrl   = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
	tun2socksDownloadUrl = "https://github.com/heiher/hev-socks5-tunnel/releases/latest/download/hev-socks5-tunnel-linux-arm64"
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
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [core|tun2socks|geodata|subscribe]").WithPrefix("update").WithPathObj(*this)
	}
	return nil
}

// updateCore update core, support xray, v2fly, sagernet
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
	coreZipPath := path.Join(builds.Config.XrayHelper.DataDir, "core.zip")
	switch builds.Config.XrayHelper.CoreType {
	case "xray":
		if err := common.DownloadFile(coreZipPath, xrayCoreDownloadUrl); err != nil {
			return err
		}
	default:
		return errors.New("unknown core type " + builds.Config.XrayHelper.CoreType).WithPrefix("update")
	}
	// update core need stop core service first
	if len(getServicePid()) > 0 {
		log.HandleInfo("update: detect core is running, stop it")
		stopService()
		serviceRunFlag = true
		_ = os.Remove(builds.Config.XrayHelper.CorePath)
	}
	zipReader, err := zip.OpenReader(coreZipPath)
	if err != nil {
		return errors.New("open core.zip failed, ", err).WithPrefix("update")
	}
	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
		_ = os.Remove(coreZipPath)
	}(zipReader)
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "ray") {
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
	if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geoip.dat"), geoipDownloadUrl); err != nil {
		return err
	}
	if err := common.DownloadFile(path.Join(builds.Config.XrayHelper.DataDir, "geosite.dat"), geositeDownloadUrl); err != nil {
		return err
	}
	return nil
}

// updateSubscribe update subscribe
func updateSubscribe() error {
	if err := os.MkdirAll(builds.Config.XrayHelper.DataDir, 0644); err != nil {
		return errors.New("create DataDir failed, ", err).WithPrefix("update")
	}
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
	return nil
}
