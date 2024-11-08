package routes

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"encoding/json"
)

const tagDnsrule = "dnsrule"

var dnsrule serial.OrderedArray

// loadDnsrule load current dns rules from core config
func loadDnsrule() {
	if len(dnsrule) > 0 {
		return
	}
	read := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagDnsrule)
		}
		if d, ok := jsonMap.Get("dns"); ok {
			dnsMap := d.Value.(serial.OrderedMap)
			if rules, ok := dnsMap.Get("rules"); ok {
				dnsrule = rules.Value.(serial.OrderedArray)
				return false, nil, nil
			}
		}
		return false, nil, e.New("cannot find dns rules from your config").WithPrefix(tagDnsrule)
	}
	_ = common.HandleCoreConfDir(read)
}

// AddDnsrule add a dns rule
func AddDnsrule(r *serial.OrderedMap) bool {
	loadDnsrule()
	dnsrule = append(dnsrule, *r)
	return true
}

// DeleteDnsrule delete a dns rule
func DeleteDnsrule(index int) bool {
	loadDnsrule()
	if index >= 0 && index < len(dnsrule) {
		dnsrule = append(dnsrule[:index], dnsrule[index+1:]...)
		return true
	}
	return false
}

// SetDnsrule replace a dns rule
func SetDnsrule(index int, r *serial.OrderedMap) bool {
	loadDnsrule()
	if index >= 0 && index < len(dnsrule) {
		dnsrule[index] = *r
		return true
	}
	return false
}

// ExchangeDnsrule exchange two dns rules' order
func ExchangeDnsrule(a int, b int) bool {
	loadDnsrule()
	if a >= 0 && a < len(dnsrule) && b >= 0 && b < len(dnsrule) {
		tmp := dnsrule[a]
		dnsrule[a] = dnsrule[b]
		dnsrule[b] = tmp
		return true
	}
	return false
}

// GetDnsrule get dns rules
func GetDnsrule() serial.OrderedArray {
	loadDnsrule()
	return dnsrule
}

// ApplyDnsrule sync dns rules to core config
func ApplyDnsrule() error {
	replace := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagDnsrule)
		}
		if dnsMap, ok := jsonMap.Get("dns"); ok {
			d := dnsMap.Value.(serial.OrderedMap)
			if _, ok := d.Get("rules"); ok {
				d.Set("rules", dnsrule)
			}
			// replace
			jsonMap.Set("dns", d)
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return false, nil, e.New("marshal config json failed, ", err).WithPrefix(tagDnsrule)
			}
			return true, marshal, nil
		}
		return false, nil, e.New("cannot found dns rules from your config").WithPrefix(tagDnsrule)
	}
	return common.HandleCoreConfDir(replace)
}
