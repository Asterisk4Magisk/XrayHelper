package routes

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"encoding/json"
	"os"
	"path"
	"strings"
)

const tagRuleset = "ruleset"

var ruleset serial.OrderedArray

// loadRuleset load current ruleset from core config
func loadRuleset() {
	if len(ruleset) > 0 {
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
		case "sing-box":
			if route, ok := jsonMap.Get("route"); ok {
				routeMap := route.Value.(serial.OrderedMap)
				if ruleSet, ok := routeMap.Get("rule_set"); ok {
					ruleset = ruleSet.Value.(serial.OrderedArray)
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

// AddRuleset add a ruleset
func AddRuleset(r *serial.OrderedMap) bool {
	loadRuleset()
	ruleset = append(ruleset, *r)
	return true
}

// DeleteRuleset delete a ruleset
func DeleteRuleset(index int) bool {
	loadRuleset()
	if index >= 0 && index < len(ruleset) {
		ruleset = append(ruleset[:index], ruleset[index+1:]...)
		return true
	}
	return false
}

// SetRuleset replace a ruleset
func SetRuleset(index int, r *serial.OrderedMap) bool {
	loadRuleset()
	if index >= 0 && index < len(ruleset) {
		ruleset[index] = *r
		return true
	}
	return false
}

// GetRuleset get the ruleset
func GetRuleset() serial.OrderedArray {
	loadRuleset()
	return ruleset
}

// ApplyRuleset sync ruleset to core config
func ApplyRuleset() error {
	replace := func(c []byte) ([]byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRuleset)
		}
		replaced := false
		switch builds.Config.XrayHelper.CoreType {
		case "sing-box":
			if routeMap, ok := jsonMap.Get("route"); ok {
				route := routeMap.Value.(serial.OrderedMap)
				if _, ok := route.Get("rule_set"); ok {
					route.Set("rule_set", ruleset)
				}
				// replace
				jsonMap.Set("route", route)
				replaced = true
			}
		}
		if replaced {
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return nil, e.New("marshal config json failed, ", err).WithPrefix(tagRuleset)
			}
			return marshal, nil
		} else {
			return nil, e.New("cannot found ruleset from your config").WithPrefix(tagRuleset)
		}
	}
	confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig)
	if err != nil {
		return e.New("open core config file failed, " + err.Error()).WithPrefix(tagRuleset)
	}
	if confInfo.IsDir() {
		if confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig); err == nil {
			for _, conf := range confDir {
				if !conf.IsDir() && strings.HasSuffix(conf.Name(), ".json") {
					if confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name())); err == nil {
						if confByte, err = replace(confByte); err == nil {
							if err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), confByte, 0644); err == nil {
								break
							} else {
								log.HandleDebug("write new config failed, " + err.Error())
							}
						} else {
							log.HandleDebug(err)
						}
					}
				}
			}
		}
	} else {
		if confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig); err == nil {
			if confByte, err = replace(confByte); err == nil {
				if err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig), confByte, 0644); err != nil {
					return e.New("write new config failed, " + err.Error()).WithPrefix(tagRuleset)
				}
			} else {
				return err
			}
		} else {
			return e.New("read core config file failed").WithPrefix(tagRuleset)
		}
	}
	return nil
}
