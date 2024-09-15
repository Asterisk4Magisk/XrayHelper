package trojan_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testTrojan = "trojan://asd-asfasf-asfasf@tj.com:443?mode=multi&security=reality&alpn=h2&pbk=111&fp=ios&spx=333&type=grpc&serviceName=wwwssss&sni=baidu.com&sid=222#tj"

func TestTrojan(t *testing.T) {
	trojanShareUrl, err := shareurls.Parse(testTrojan)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(trojanShareUrl.GetNodeInfoStr())
	tag, err := trojanShareUrl.ToOutboundWithTag("xray", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(string(indent))
}
