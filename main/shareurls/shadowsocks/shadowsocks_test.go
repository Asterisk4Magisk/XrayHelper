package shadowsocks_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"encoding/json"
	"github.com/tailscale/hujson"
	"os"
	"testing"
)

const testSS = "ss://YWVzLTI1Ni1nY206dGVzdHNoYWRvd3NvY2tz@0.0.0.0:65535#%E6%B5%8B%E8%AF%95SS"

func TestShadowsocks(t *testing.T) {
	// read origin json
	configFile, err := os.ReadFile("/XrayHelper/proxy.json")
	if err != nil {
		t.Error(err)
	}
	// standardize origin json (remove comment)
	standardize, err := hujson.Standardize(configFile)
	if err != nil {
		t.Error(err)
	}
	// unmarshal
	var jsonValue interface{}
	err = json.Unmarshal(standardize, &jsonValue)
	if err != nil {
		t.Error(err)
	}
	// assert json to map
	jsonMap, ok := jsonValue.(map[string]interface{})
	if !ok {
		t.Error("convert to json map error")
	}
	outbounds, ok := jsonMap["outbounds"]
	if !ok {
		t.Error("cannot find outbounds")
	}
	// assert outbounds
	outboundsMap, ok := outbounds.([]interface{})
	if !ok {
		t.Error("outbounds is invalid")
	}
	for i, outbound := range outboundsMap {
		outboundMap, ok := outbound.(map[string]interface{})
		if !ok {
			continue
		}
		tag, ok := outboundMap["tag"].(string)
		if !ok {
			continue
		}
		if tag == "proxy" {
			log.HandleInfo(outbound)
			// new a shareUrl with shadowsocks url
			shareUrl, err := shareurls.ParseShareUrl(testSS)
			if err != nil {
				t.Error(err)
			}
			log.HandleInfo(shareUrl.GetNodeInfo())
			// replace
			outbound, err = shareUrl.ToOutoundWithTag("xray", "proxy")
			if err != nil {
				t.Error(err)
			}
			outboundsMap[i] = outbound
			jsonMap["outbounds"] = outboundsMap
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				t.Error(err)
			}
			log.HandleInfo(string(marshal))
		}
	}
}
