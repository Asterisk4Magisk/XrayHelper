package routes

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"encoding/json"
	"os"
	"path"
	"strings"
)

const tagRule = "rule"

var rule serial.OrderedArray

func loadRule() {
	if len(rule) > 0 {
		return
	}
	read := func(c []byte) bool {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			log.HandleDebug("json unmarshal failed, " + err.Error())
			return false
		}
		switch builds.Config.XrayHelper.CoreType {
		case "xray":
			if routing, ok := jsonMap.Get("routing"); ok {
				routingMap := routing.Value.(serial.OrderedMap)
				if rules, ok := routingMap.Get("rules"); ok {
					rule = rules.Value.(serial.OrderedArray)
					return true
				}
			}
		case "sing-box":
			if route, ok := jsonMap.Get("route"); ok {
				routeMap := route.Value.(serial.OrderedMap)
				if rules, ok := routeMap.Get("rules"); ok {
					rule = rules.Value.(serial.OrderedArray)
					return true
				}
			}
		}
		return false
	}
	confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig)
	if err != nil {
		log.HandleDebug("open core config file failed, " + err.Error())
		return
	}
	if confInfo.IsDir() {
		if confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig); err == nil {
			for _, conf := range confDir {
				if !conf.IsDir() && strings.HasSuffix(conf.Name(), ".json") {
					if confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name())); err == nil {
						if read(confByte) {
							break
						}
					}
				}
			}
		}
	} else {
		if confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig); err == nil {
			read(confByte)
		}
	}
}

func GetRule() serial.OrderedArray {
	loadRule()
	return rule
}
