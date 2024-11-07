package routes

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/switches"
	"encoding/json"
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
	read := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRule)
		}
		switch builds.Config.XrayHelper.CoreType {
		case "xray":
			if routing, ok := jsonMap.Get("routing"); ok {
				routingMap := routing.Value.(serial.OrderedMap)
				if rules, ok := routingMap.Get("rules"); ok {
					rule = rules.Value.(serial.OrderedArray)
					return false, nil, nil
				}
			}
		case "sing-box":
			if route, ok := jsonMap.Get("route"); ok {
				routeMap := route.Value.(serial.OrderedMap)
				if rules, ok := routeMap.Get("rules"); ok {
					rule = rules.Value.(serial.OrderedArray)
					return false, nil, nil
				}
			}
		}
		return false, nil, e.New("cannot find rule from your config").WithPrefix(tagRule)
	}
	_ = common.HandleCoreConfDir(read)
}

// AddRule add a rule
func AddRule(r *serial.OrderedMap) bool {
	loadRule()
	rule = append(rule, *r)
	return true
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

// ExchangeRule exchange two rules' order
func ExchangeRule(a int, b int) bool {
	loadRule()
	if a >= 0 && a < len(rule) && b >= 0 && b < len(rule) {
		tmp := rule[a]
		rule[a] = rule[b]
		rule[b] = tmp
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
	replace := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRule)
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
				return false, nil, e.New("marshal config json failed, ", err).WithPrefix(tagRule)
			}
			return true, marshal, nil
		} else {
			return false, nil, e.New("cannot found rules from your config").WithPrefix(tagRule)
		}
	}
	return common.HandleCoreConfDir(replace)
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
	replace := func(c []byte) (bool, []byte, error) {
		s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType)
		if err != nil {
			return false, nil, err
		}
		var jsonMap serial.OrderedMap
		err = json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRule)
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
				return false, nil, e.New("marshal config json failed, ", err).WithPrefix(tagRule)
			}
			return true, marshal, nil
		}
		return false, nil, e.New("cannot found outbounds from your config").WithPrefix(tagRule)
	}
	return common.HandleCoreConfDir(replace)
}
