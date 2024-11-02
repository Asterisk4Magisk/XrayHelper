package routes

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/switches"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
)

const tagRule = "rule"

var rule serial.OrderedArray

// loadRule load current rules from core config
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

// AddRule add a rule
func AddRule(index int, r *serial.OrderedMap) bool {
	loadRule()
	if index >= 0 && index <= len(rule) {
		rule = append(rule[:index+1], rule[index:]...)
		rule[index] = *r
		return true
	}
	return false
}

// DeleteRule delete a rule
func DeleteRule(index int) bool {
	loadRule()
	if index >= 0 && index < len(rule) {
		rule = append(rule[:index], rule[index+1:]...)
		return true
	}
	return false
}

// SetRule replace a rule
func SetRule(index int, r *serial.OrderedMap) bool {
	loadRule()
	if index >= 0 && index < len(rule) {
		rule[index] = *r
		return true
	}
	return false
}

// GetRule get rules
func GetRule() serial.OrderedArray {
	loadRule()
	return rule
}

// ApplyRule sync rules to core config
func ApplyRule() error {
	if err := replaceOutbounds(); err != nil {
		return err
	}
	replace := func(c []byte) ([]byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRule)
		}
		replaced := false
		switch builds.Config.XrayHelper.CoreType {
		case "xray":
			if routingMap, ok := jsonMap.Get("routing"); ok {
				routing := routingMap.Value.(serial.OrderedMap)
				if _, ok := routing.Get("rules"); ok {
					routing.Set("rules", rule)
				}
				// replace
				jsonMap.Set("routing", routing)
				replaced = true
			}
		case "sing-box":
			if routeMap, ok := jsonMap.Get("route"); ok {
				route := routeMap.Value.(serial.OrderedMap)
				if _, ok := route.Get("rules"); ok {
					route.Set("rules", rule)
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
				return nil, e.New("marshal config json failed, ", err).WithPrefix(tagRule)
			}
			return marshal, nil
		} else {
			return nil, e.New("cannot found rules from your config").WithPrefix(tagRule)
		}
	}
	confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig)
	if err != nil {
		return e.New("open core config file failed, " + err.Error()).WithPrefix(tagRule)
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
					return e.New("write new config failed, " + err.Error()).WithPrefix(tagRule)
				}
			} else {
				return err
			}
		} else {
			return e.New("read core config file failed").WithPrefix(tagRule)
		}
	}
	return nil
}

func replaceOutbounds() error {
	getOutBoundTags := func() (tags []string) {
		var tagName = "outboundTag"
		switch builds.Config.XrayHelper.CoreType {
		case "xray":
			tagName = "outboundTag"
		case "sing-box":
			tagName = "outbound"
		}
		for _, r := range rule {
			ruleMap := r.(serial.OrderedMap)
			if tag, ok := ruleMap.Get(tagName); ok {
				tags = append(tags, tag.Value.(string))
			}
		}
		return
	}
	replace := func(c []byte) ([]byte, error) {
		s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType)
		if err != nil {
			return nil, err
		}
		var jsonMap serial.OrderedMap
		err = json.Unmarshal(c, &jsonMap)
		if err != nil {
			return nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRule)
		}
		if outbounds, ok := jsonMap.Get("outbounds"); ok {
			outboundsArray := outbounds.Value.(serial.OrderedArray)
			for i := 0; i < len(outboundsArray); i++ {
				outboundMap := outboundsArray[i].(serial.OrderedMap)
				if tag, ok := outboundMap.Get("tag"); ok {
					if strings.HasPrefix(tag.Value.(string), "xrayhelper") {
						outboundsArray = append(outboundsArray[:i], outboundsArray[i+1:]...)
						i--
					}
				}
			}
			// collect
			var subscribe, custom []int
			for _, tag := range getOutBoundTags() {
				if strings.HasPrefix(tag, "xrayhelper-") {
					if index, err := strconv.Atoi(strings.TrimPrefix(tag, "xrayhelper-")); err == nil {
						subscribe = append(subscribe, index)
					}
				} else if strings.HasPrefix(tag, "xrayhelpercustom-") {
					if index, err := strconv.Atoi(strings.TrimPrefix(tag, "xrayhelpercustom-")); err == nil {
						custom = append(custom, index)
					}
				}
			}
			for _, i := range subscribe {
				tag := "xrayhelper-" + strconv.Itoa(i)
				shareurl := s.Choose(false, i).(shareurls.ShareUrl)
				if o, err := shareurl.ToOutboundWithTag(builds.Config.XrayHelper.CoreType, tag); err == nil {
					outboundsArray = append(outboundsArray, o)
				}
			}
			s.Clear()
			for _, i := range custom {
				tag := "xrayhelpercustom-" + strconv.Itoa(i)
				shareurl := s.Choose(true, i).(shareurls.ShareUrl)
				if o, err := shareurl.ToOutboundWithTag(builds.Config.XrayHelper.CoreType, tag); err == nil {
					outboundsArray = append(outboundsArray, o)
				}
			}
			// replace
			jsonMap.Set("outbounds", outboundsArray)
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return nil, e.New("marshal config json failed, ", err).WithPrefix(tagRule)
			}
			return marshal, nil
		}
		return nil, e.New("cannot found outbounds from your config").WithPrefix(tagRule)
	}
	confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig)
	if err != nil {
		return e.New("open core config file failed, " + err.Error()).WithPrefix(tagRule)
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
					return e.New("write new config failed, " + err.Error()).WithPrefix(tagRule)
				}
			} else {
				return err
			}
		} else {
			return e.New("read core config file failed").WithPrefix(tagRule)
		}
	}
	return nil
}
