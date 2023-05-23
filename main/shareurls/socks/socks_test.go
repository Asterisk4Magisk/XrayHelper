package socks_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"testing"
)

const testSocks = "socks://cXdlOmFzZA==@socks5.com:443#%E6%B5%8B%E8%AF%95SOCKS"

func TestSocks(t *testing.T) {
	socksShareUrl, err := shareurls.ParseShareUrl(testSocks)
	if err != nil {
		t.Error(err)
	}
	log.HandleInfo(socksShareUrl.GetNodeInfo())
}
