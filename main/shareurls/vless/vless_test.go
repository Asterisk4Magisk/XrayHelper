package vless_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testVLESS = "vless://6666-66666666-666666@1.com:443?path=%2Fcccc&security=tls&encryption=none&alpn=h2,http/1.1&host=2.com&fp=firefox&type=http&flow=xtls-rprx-vision&sni=3.com#%E6%B5%8B%E8%AF%95%E8%8A%82%E7%82%B9"

func TestVLESS(t *testing.T) {
	vlessShareUrl, err := shareurls.Parse(testVLESS)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(vlessShareUrl.GetNodeInfo())
}
