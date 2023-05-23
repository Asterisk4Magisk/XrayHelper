package trojan_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testTrojan = "trojan://asd-asfasf-asfasf@tj.com:443?mode=multi&security=reality&alpn=h2&pbk=111&fp=ios&spx=333&type=grpc&serviceName=wwwssss&sni=baidu.com&sid=222#tj"

func TestTrojan(t *testing.T) {
	trojanShareUrl, err := shareurls.ParseShareUrl(testTrojan)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(trojanShareUrl.GetNodeInfo())
}
