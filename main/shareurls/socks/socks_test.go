package socks_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testSocks = "socks://cXdlOmFzZA==@socks5.com:443#%E6%B5%8B%E8%AF%95SOCKS"

func TestSocks(t *testing.T) {
	socksShareUrl, err := shareurls.Parse(testSocks)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(socksShareUrl.GetNodeInfo())
	tag, err := socksShareUrl.ToOutboundWithTag("xray", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(string(indent))
}
