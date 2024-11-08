package routes

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"encoding/json"
)

const tagRuleset = "ruleset"

var ruleset serial.OrderedArray

// loadRuleset load current ruleset from core config
func loadRuleset() {
	if len(ruleset) > 0 {
		return
	}
	read := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRuleset)
		}
		if route, ok := jsonMap.Get("route"); ok {
			routeMap := route.Value.(serial.OrderedMap)
			if ruleSet, ok := routeMap.Get("rule_set"); ok {
				ruleset = ruleSet.Value.(serial.OrderedArray)
				return false, nil, nil
			}
		}
		return false, nil, e.New("cannot find rule_set from your config").WithPrefix(tagRuleset)
	}
	_ = common.HandleCoreConfDir(read)
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
	replace := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagRuleset)
		}
		if routeMap, ok := jsonMap.Get("route"); ok {
			route := routeMap.Value.(serial.OrderedMap)
			if _, ok := route.Get("rule_set"); ok {
				route.Set("rule_set", ruleset)
			}
			// replace
			jsonMap.Set("route", route)
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return false, nil, e.New("marshal config json failed, ", err).WithPrefix(tagRuleset)
			}
			return true, marshal, nil
		}
		return false, nil, e.New("cannot found ruleset from your config").WithPrefix(tagRuleset)
	}
	return common.HandleCoreConfDir(replace)
}
