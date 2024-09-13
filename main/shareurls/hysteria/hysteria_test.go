package hysteria_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testHysteria = "hysteria://hysteria.network:443?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks"

func TestHysteria(t *testing.T) {
	hysteriaShareUrl, err := shareurls.Parse(testHysteria)
	if err != nil {
		t.Error(err)
	}
	tag, err := hysteriaShareUrl.ToOutboundWithTag("sing-box", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(hysteriaShareUrl.GetNodeInfoStr())
	fmt.Println(string(indent))
}
