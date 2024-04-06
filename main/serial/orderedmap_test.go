package serial

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	yamlStr := []byte(`xrayHelper:
    coreType: sing-box
    corePath: /data/adb/xray/bin/sing-box
    coreConfig: /data/adb/xray/singconfs/
    dataDir: /data/adb/xray/data/
    runDir: /data/adb/xray/run/
    proxyTag: proxy
    subList:
        - obj: 111
          xxx: 222
          yyy:
              zzz: 333
              uuu: 444
          ttt: 000
proxy:
    method: tproxy
    tproxyPort: 65535
    socksPort: 65534
    tunDevice: xtun
    enableIPv6: false
    autoDNSStrategy: false
    mode: blacklist
    apList:
        - rndis+
        - wlan+
    intraList:
        - 10.10.10.0/24
clash:
    dnsPort: 65533
    template: /data/adb/xray/mihomoconfs/template.yaml`)
	var yamlMap OrderedMap
	err := yaml.Unmarshal(yamlStr, &yamlMap)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	marshal, _ := yaml.Marshal(yamlMap)
	fmt.Print(string(marshal))
}
