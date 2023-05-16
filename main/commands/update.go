package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/utils"
	"archive/zip"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	xrayCoreUrl     = "https://github.com/XTLS/Xray-core/releases/latest/download/Xray-android-arm64-v8a.zip"
	v2flyCoreUrl    = "https://github.com/v2fly/v2ray-core/releases/latest/download/v2ray-android-arm64-v8a.zip"
	sagernetCoreUrl = "https://github.com/SagerNet/v2ray-core/releases/latest/download/v2ray-android-arm64-v8a.zip"
	geoipUrl        = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
	geositeUrl      = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
)

type UpdateCommand struct{}

func (this *UpdateCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify operation, available operation [core|geodata]").WithPrefix("update").WithPathObj(*this)
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
	case "geodata":
		log.HandleInfo("update: updating geodata")
		if err := updateGeodata(); err != nil {
			return err
		}
		log.HandleInfo("update: update success")
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [core|geodata]").WithPrefix("update").WithPathObj(*this)
	}
	return nil
}

// updateCore update core, support xray, v2fly, sagernet
func updateCore() error {
	if runtime.GOARCH != "arm64" {
		return errors.New("this feature only support arm64 device").WithPrefix("update")
	}
	serviceRunFlag := false
	coreZipPath := path.Join(builds.Config.XrayHelper.RunDir, "core.zip")
	switch builds.Config.XrayHelper.CoreType {
	case "xray":
		if err := utils.DownloadFile(coreZipPath, xrayCoreUrl); err != nil {
			return err
		}
	case "v2fly":
		if err := utils.DownloadFile(coreZipPath, v2flyCoreUrl); err != nil {
			return err
		}
	case "sagernet":
		if err := utils.DownloadFile(coreZipPath, sagernetCoreUrl); err != nil {
			return err
		}
	default:
		return errors.New("unknown core type " + builds.Config.XrayHelper.CoreType).WithPrefix("update")
	}
	// update core need stop core service first
	if len(getServicePid()) > 0 {
		stopService()
		serviceRunFlag = true
		_ = os.Remove(builds.Config.XrayHelper.CorePath)
	}
	zipReader, err := zip.OpenReader(coreZipPath)
	if err != nil {
		return errors.New("open core.zip failed ,", err).WithPrefix("update")
	}
	defer func(zipReader *zip.ReadCloser) {
		_ = zipReader.Close()
		_ = os.Remove(coreZipPath)
	}(zipReader)
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "ray") {
			savePath := path.Join(path.Dir(builds.Config.XrayHelper.CorePath), file.Name)
			fileReader, err := file.Open()
			if err != nil {
				return errors.New("cannot get file reader "+file.Name+", ", err).WithPrefix("update")
			}
			saveFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
			if err != nil {
				return errors.New("cannot open file "+savePath+", ", err).WithPrefix("update")
			}
			_, err = io.Copy(saveFile, fileReader)
			if err != nil {
				return errors.New("save file "+savePath+" failed, ", err).WithPrefix("net")
			}
			_ = saveFile.Close()
			_ = fileReader.Close()
			break
		}
	}
	if serviceRunFlag {
		_ = startService()
	}
	return nil
}

// updateGeodata update geodata
func updateGeodata() error {
	if err := utils.DownloadFile(path.Join(builds.Config.XrayHelper.GeodataDir, "geoip.dat"), geoipUrl); err != nil {
		return err
	}
	if err := utils.DownloadFile(path.Join(builds.Config.XrayHelper.GeodataDir, "geosite.dat"), geositeUrl); err != nil {
		return err
	}
	return nil
}
