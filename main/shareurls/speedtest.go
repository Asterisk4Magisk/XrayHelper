package shareurls

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"encoding/json"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	tagSpeedtest = "speedtest"
	testFileName = "test.json"
	testPort     = "65500"
	testUrl      = "https://www.google.com/generate_204"
)

func RealPing(coreType string, url ShareUrl) (result int) {
	result = -1
	configPath := path.Join(builds.Config.XrayHelper.RunDir, testFileName)
	// start test service
	service, err := startTestService(coreType, url, configPath)
	if err != nil {
		log.HandleDebug(err)
		return
	}
	// check service port
	listenFlag := false
	for i := 0; i < *builds.CoreStartTimeout; i++ {
		time.Sleep(1 * time.Second)
		if common.CheckLocalPort(strconv.Itoa(service.Pid()), testPort, false) {
			listenFlag = true
			break
		}
	}
	if !listenFlag {
		log.HandleDebug("service not listen for RealPing")
		return
	}
	// set socks5 proxy
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+testPort, nil, proxy.Direct)
	if err != nil {
		log.HandleDebug("set socks5 proxy: " + err.Error())
		return
	}
	client := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	// start test
	request, _ := http.NewRequest("GET", testUrl, nil)
	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		log.HandleDebug("request google_204: " + err.Error())
		return
	}
	// defer stop test service
	defer func(Body io.ReadCloser, Service common.External) {
		_ = Body.Close()
		stopTest(Service, configPath)
	}(response.Body, service)
	// get result
	if response.StatusCode != 204 {
		log.HandleDebug("request google_204 get " + strconv.Itoa(response.StatusCode))
		return
	}
	result = int(time.Since(start).Milliseconds())
	return
}

func startTestService(coreType string, url ShareUrl, configPath string) (common.External, error) {
	var service common.External
	switch coreType {
	case "xray":
		if err := genXrayTestConfig(url, configPath); err != nil {
			return nil, err
		}
		service = common.NewExternal(0, nil, nil, builds.Config.XrayHelper.CorePath, "run", "-c", configPath)
	default:
		return nil, e.New("not a supported coreType " + coreType).WithPrefix(tagSpeedtest)
	}
	service.SetUidGid("0", common.CoreGid)
	service.Start()
	return service, nil
}

func genXrayTestConfig(url ShareUrl, configPath string) error {
	var config serial.OrderedMap
	// add dns
	var dnsObj serial.OrderedMap
	var dnsServersArr serial.OrderedArray
	dnsServersArr = append(dnsServersArr, "1.1.1.1", "223.5.5.5")
	dnsObj.Set("servers", dnsServersArr)
	config.Set("dns", dnsObj)
	// add inbounds
	var inboundsArr serial.OrderedArray
	var socksObj serial.OrderedMap
	socksObj.Set("port", testPort)
	socksObj.Set("protocol", "socks")
	inboundsArr = append(inboundsArr, socksObj)
	config.Set("inbounds", inboundsArr)
	// add outbounds
	var outboundsArr serial.OrderedArray
	outbound, err := url.ToOutboundWithTag("xray", "test")
	if err != nil {
		return err
	}
	outboundsArr = append(outboundsArr, outbound)
	config.Set("outbounds", outboundsArr)
	// save test config
	marshal, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return e.New("marshal xray test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	if err := os.WriteFile(configPath, marshal, 0644); err != nil {
		return e.New("write xray test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	return nil
}

func stopTest(service common.External, configPath string) {
	_ = service.Kill()
	//_ = os.Remove(configPath)
}
