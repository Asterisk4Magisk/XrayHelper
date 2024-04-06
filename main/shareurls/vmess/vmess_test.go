package vmess_test

import (
	"XrayHelper/main/shareurls"
	"encoding/json"
	"fmt"
	"testing"
)

const testVmess = "vmess://eyJhZGQiOiIzMjEuY29tIiwiYWlkIjoiMiIsImFscG4iOiJoMiIsImZwIjoiZWRnZSIsImhvc3QiOiIiLCJpZCI6IjY2NjYtNjY2Ni02NjY2IiwibmV0IjoidGNwIiwicGF0aCI6IiIsInBvcnQiOiI0NDMiLCJwcyI6IjMyMSIsInNjeSI6ImFlcy0xMjgtZ2NtIiwic25pIjoiIiwidGxzIjoidGxzIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiJ9"

func TestVmess(t *testing.T) {
	vmessShareUrl, err := shareurls.Parse(testVmess)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vmessShareUrl.GetNodeInfo())
	tag, err := vmessShareUrl.ToOutboundWithTag("xray", "proxy")
	indent, err := json.MarshalIndent(tag, "", "    ")
	fmt.Println(string(indent))
}
