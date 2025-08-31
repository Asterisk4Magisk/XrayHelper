package vless_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testVLESS = "vless://id@host:444?type=xhttp&mode=auto&host=host.com&path=%2Fpath&extra=%7B%0A%20%20%20%20%20%20%20%20%22headers%22%3A%20%7B%0A%20%20%20%20%20%20%20%20%20%20%20%20%22kkk%22%3A%20%22vvv%22%0A%20%20%20%20%20%20%20%20%7D%2C%0A%20%20%20%20%20%20%20%20%22xPaddingBytes%22%3A%20%22100-1000%22%2C%0A%20%20%20%20%20%20%20%20%22noGRPCHeader%22%3A%20false%2C%20%0A%20%20%20%20%20%20%20%20%22noSSEHeader%22%3A%20false%2C%20%0A%20%20%20%20%20%20%20%20%22scMaxEachPostBytes%22%3A%201000000%2C%20%0A%20%20%20%20%20%20%20%20%22scMinPostsIntervalMs%22%3A%2030%2C%20%0A%20%20%20%20%20%20%20%20%22scMaxBufferedPosts%22%3A%2030%2C%20%0A%20%20%20%20%20%20%20%20%22scStreamUpServerSecs%22%3A%20%2220-80%22%2C%20%0A%20%20%20%20%20%20%20%20%22xmux%22%3A%20%7B%20%0A%20%20%20%20%20%20%20%20%20%20%20%20%22maxConcurrency%22%3A%20%2216-32%22%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22maxConnections%22%3A%200%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22cMaxReuseTimes%22%3A%200%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22hMaxRequestTimes%22%3A%20%22600-900%22%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22hMaxReusableSecs%22%3A%20%221800-3000%22%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22hKeepAlivePeriod%22%3A%200%0A%20%20%20%20%20%20%20%20%7D%2C%0A%20%20%20%20%20%20%20%20%22downloadSettings%22%3A%20%7B%20%0A%20%20%20%20%20%20%20%20%20%20%20%20%22address%22%3A%20%22%22%2C%20%0A%20%20%20%20%20%20%20%20%20%20%20%20%22port%22%3A%20443%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22network%22%3A%20%22xhttp%22%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22security%22%3A%20%22tls%22%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22xhttpSettings%22%3A%20%7B%0A%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%20%22path%22%3A%20%22%2Fyourpath%22%0A%20%20%20%20%20%20%20%20%20%20%20%20%7D%2C%0A%20%20%20%20%20%20%20%20%20%20%20%20%22sockopt%22%3A%20%7B%7D%20%0A%20%20%20%20%20%20%20%20%7D%0A%20%20%20%20%7D&#remarks"

func TestVLESS(t *testing.T) {
	vlessShareUrl, err := shareurls.Parse(testVLESS)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vlessShareUrl.GetNodeInfoStr())
	tag, err := vlessShareUrl.ToOutboundWithTag("xray", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(string(indent))
}
