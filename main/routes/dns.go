package routes

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"encoding/json"
)

const tagDns = "dns"

var dns serial.OrderedArray

// loadDns load current dns servers from core config
func loadDns() {
	if len(dns) > 0 {
		return
	}
	read := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagDns)
		}
		if d, ok := jsonMap.Get("dns"); ok {
			dnsMap := d.Value.(serial.OrderedMap)
			if servers, ok := dnsMap.Get("servers"); ok {
				dns = servers.Value.(serial.OrderedArray)
				return false, nil, nil
			}
		}
		return false, nil, e.New("cannot find dns servers from your config").WithPrefix(tagDns)
	}
	_ = common.HandleCoreConfDir(read)
}

// AddDns add a dns
func AddDns[T any](r *T) bool {
	loadDns()
	dns = append(dns, *r)
	return true
}

// DeleteDns delete a dns
func DeleteDns(index int) bool {
	loadDns()
	if index >= 0 && index < len(dns) {
		dns = append(dns[:index], dns[index+1:]...)
		return true
	}
	return false
}

// SetDns replace a dns
func SetDns[T any](index int, r *T) bool {
	loadDns()
	if index >= 0 && index < len(dns) {
		dns[index] = *r
		return true
	}
	return false
}

// GetDns get dns
func GetDns() serial.OrderedArray {
	loadDns()
	return dns
}

// ApplyDns sync dns servers to core config
func ApplyDns() error {
	replace := func(c []byte) (bool, []byte, error) {
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(c, &jsonMap)
		if err != nil {
			return false, nil, e.New("json unmarshal failed, " + err.Error()).WithPrefix(tagDns)
		}
		if dnsMap, ok := jsonMap.Get("dns"); ok {
			d := dnsMap.Value.(serial.OrderedMap)
			if _, ok := d.Get("servers"); ok {
				d.Set("servers", dns)
			}
			// replace
			jsonMap.Set("dns", d)
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return false, nil, e.New("marshal config json failed, ", err).WithPrefix(tagDns)
			}
			return true, marshal, nil
		}
		return false, nil, e.New("cannot found dns from your config").WithPrefix(tagDns)
	}
	return common.HandleCoreConfDir(replace)
}
