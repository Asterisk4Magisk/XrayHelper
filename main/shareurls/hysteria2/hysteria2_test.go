package hysteria2_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testHysteria2 = "hysteria2://letmein:123456@example.com:443/?insecure=1&obfs=salamander&obfs-password=gawrgura&pinSHA256=deadbeef&sni=real.example.com#Remark"

func TestHysteria2(t *testing.T) {
	hysteriaShareUrl, err := shareurls.Parse(testHysteria2)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(hysteriaShareUrl.GetNodeInfo())
}
