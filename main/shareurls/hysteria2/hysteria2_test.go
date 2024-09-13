package hysteria2_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testHysteria2 = "hysteria2://letmein:123456@example.com:443/?insecure=1&obfs=salamander&obfs-password=gawrgura&pinSHA256=deadbeef&sni=real.example.com#Remark"

func TestHysteria2(t *testing.T) {
	hysteriaShareUrl, err := shareurls.Parse(testHysteria2)
	if err != nil {
		t.Error(err)
	}
	tag, err := hysteriaShareUrl.ToOutboundWithTag("sing-box", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(hysteriaShareUrl.GetNodeInfoStr())
	fmt.Println(string(indent))
}
