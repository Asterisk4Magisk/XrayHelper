package hysteria_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testHysteria = "hysteria://hysteria.network:443?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks"

func TestSocks(t *testing.T) {
	hysteriaShareUrl, err := shareurls.Parse(testHysteria)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(hysteriaShareUrl.GetNodeInfo())
}
