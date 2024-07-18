package shadowsocks_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testSS = "ss://YWVzLTI1Ni1nY206dGVzdHNoYWRvd3NvY2tz@0.0.0.0:65535#%E6%B5%8B%E8%AF%95SS"

func TestShadowsocks(t *testing.T) {
	ssShareUrl, err := shareurls.Parse(testSS)
	if err != nil {
		t.Error(err)
	}
	tag, err := ssShareUrl.ToOutboundWithTag("xray", "proxy")
	out, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(string(out))
}
