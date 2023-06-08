package shadowsocksr_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testShadowsocksR = "ssr://NDUuMzIuMTMxLjExMTo4OTg5Om9yaWdpbjphZXMtMjU2LWNmYjpwbGFpbjpiM0JsYm5ObGMyRnRaUS8_cmVtYXJrcz1VMU5TVkU5UFRGOU9iMlJsT3VlLWp1V2J2U0RsaXFEbGlLbm5wb19sc0x6a3Vwcmx0NTdsbktQa3ZaWGxvWjVEYUc5dmNHSG1sYkRtamE3a3VLM2x2NE0mZ3JvdXA9VjFkWExsTlRVbFJQVDB3dVEwOU4"

func TestShadowsocksR(t *testing.T) {
	shadowsocksRShareUrl, err := shareurls.Parse(testShadowsocksR)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(shadowsocksRShareUrl.GetNodeInfo())
}
