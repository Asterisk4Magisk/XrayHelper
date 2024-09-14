package shareurls

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"context"
	"encoding/json"
	"golang.org/x/net/proxy"
	"io"
	"net"
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
	defer stopTestService(service, configPath)
	// check service port
	listenFlag := false
	for i := 0; i < *builds.CoreStartTimeout; i++ {
		time.Sleep(1 * time.Second)
		if common.CheckLocalPort(strconv.Itoa(service.Pid()), testPort, false) ||
			common.CheckLocalPort(strconv.Itoa(service.Pid()), testPort, true) {
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
	// drop first result
	startTest(dialer)
	result = startTest(dialer)
	return
}

func startTest(dialer proxy.Dialer) (result int) {
	result = -1
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: transport}
	// start test
	request, _ := http.NewRequest("GET", testUrl, nil)
	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		log.HandleDebug("request google_204: " + err.Error())
		return
	}
	// defer close body
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
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
		serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "test.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
		if err != nil {
			return nil, e.New("open core test log file failed, ", err).WithPrefix(tagSpeedtest)
		}
		service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", configPath)
	case "sing-box":
		if err := genSingboxTestConfig(url, configPath); err != nil {
			return nil, err
		}
		serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "test.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
		if err != nil {
			return nil, e.New("open core test log file failed, ", err).WithPrefix(tagSpeedtest)
		}
		service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", configPath, "--disable-color")
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
	dnsServersArr = append(dnsServersArr, "223.5.5.5")
	dnsObj.Set("servers", dnsServersArr)
	config.Set("dns", dnsObj)
	// add inbounds
	var inboundsArr serial.OrderedArray
	var socksObj serial.OrderedMap
	socksObj.Set("tag", "socks-in")
	port, _ := strconv.Atoi(testPort)
	socksObj.Set("port", port)
	socksObj.Set("protocol", "socks")

	var sniffingObj serial.OrderedMap
	sniffingObj.Set("enabled", true)
	var destOverrideArr serial.OrderedArray
	destOverrideArr = append(destOverrideArr, "http", "tls", "quic")
	sniffingObj.Set("destOverride", destOverrideArr)

	socksObj.Set("sniffing", sniffingObj)
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

func genSingboxTestConfig(url ShareUrl, configPath string) error {
	var config serial.OrderedMap
	// add dns
	var dnsObj serial.OrderedMap
	var dnsServersArr serial.OrderedArray
	var dnsServerObj serial.OrderedMap
	dnsServerObj.Set("address", "223.5.5.5")
	dnsServersArr = append(dnsServersArr, dnsServerObj)
	dnsObj.Set("servers", dnsServersArr)
	config.Set("dns", dnsObj)
	// add inbound
	var inboundsArr serial.OrderedArray
	var socksObj serial.OrderedMap
	socksObj.Set("tag", "socks-in")
	socksObj.Set("listen", "::")
	port, _ := strconv.Atoi(testPort)
	socksObj.Set("listen_port", port)
	socksObj.Set("type", "socks")
	socksObj.Set("sniff", true)
	socksObj.Set("sniff_override_destination", true)
	inboundsArr = append(inboundsArr, socksObj)
	config.Set("inbounds", inboundsArr)
	// add outbounds
	var outboundsArr serial.OrderedArray
	outbound, err := url.ToOutboundWithTag("sing-box", "test")
	if err != nil {
		return err
	}
	outboundsArr = append(outboundsArr, outbound)
	config.Set("outbounds", outboundsArr)
	// save test config
	marshal, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return e.New("marshal sing-box test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	if err := os.WriteFile(configPath, marshal, 0644); err != nil {
		return e.New("write sing-box test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	return nil
}

func stopTestService(service common.External, configPath string) {
	_ = service.Kill()
	//_ = os.Remove(configPath)
}
