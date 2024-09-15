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
	testUrl      = "https://www.google.com/generate_204"
)

type Result struct {
	Name  string
	Url   ShareUrl
	Port  int
	Value int
}

func RealPing(coreType string, results chan *Result, result *Result) {
	configPath := path.Join(builds.Config.XrayHelper.RunDir, serial.Concat("test", result.Port, ".json"))
	// start test service
	service, err := startTestService(coreType, result.Url, result.Port, configPath)
	if err != nil {
		log.HandleDebug(err)
		results <- result
		return
	}
	defer stopTestService(service, configPath)
	// check service port
	listenFlag := false
	for i := 0; i < *builds.CoreStartTimeout; i++ {
		time.Sleep(1 * time.Second)
		if common.CheckLocalPort(strconv.Itoa(service.Pid()), strconv.Itoa(result.Port), false) ||
			common.CheckLocalPort(strconv.Itoa(service.Pid()), strconv.Itoa(result.Port), true) {
			listenFlag = true
			break
		}
	}
	if !listenFlag {
		log.HandleDebug("service not listen for RealPing")
		results <- result
		return
	}
	// set socks5 proxy
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(result.Port), nil, proxy.Direct)
	if err != nil {
		log.HandleDebug("set socks5 proxy: " + err.Error())
		results <- result
		return
	}
	start := time.Now()
	for {
		result.Value = startTest(dialer)
		if time.Since(start) > 5*time.Second || result.Value > -1 {
			break
		}
	}
	results <- result
}

func startTest(dialer proxy.Dialer) (result int) {
	result = -1
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
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

func startTestService(coreType string, url ShareUrl, port int, configPath string) (common.External, error) {
	var service common.External
	switch coreType {
	case "xray":
		if err := genXrayTestConfig(url, port, configPath); err != nil {
			return nil, err
		}
		service = common.NewExternal(0, nil, nil, builds.Config.XrayHelper.CorePath, "run", "-c", configPath)
	case "sing-box":
		if err := genSingboxTestConfig(url, port, configPath); err != nil {
			return nil, err
		}
		service = common.NewExternal(0, nil, nil, builds.Config.XrayHelper.CorePath, "run", "-c", configPath, "--disable-color")
	default:
		return nil, e.New("not a supported coreType " + coreType).WithPrefix(tagSpeedtest)
	}
	service.SetUidGid("0", common.CoreGid)
	service.Start()
	return service, nil
}

func genXrayTestConfig(url ShareUrl, port int, configPath string) error {
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
	marshal, err := json.Marshal(config)
	if err != nil {
		return e.New("marshal xray test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	if err := os.WriteFile(configPath, marshal, 0644); err != nil {
		return e.New("write xray test config failed, ", err).WithPrefix(tagSpeedtest)
	}
	return nil
}

func genSingboxTestConfig(url ShareUrl, port int, configPath string) error {
	var config serial.OrderedMap
	// add dns
	var dnsObj serial.OrderedMap
	var dnsServersArr serial.OrderedArray
	var dnsServerObj serial.OrderedMap
	dnsServerObj.Set("address", "223.5.5.5")
	dnsServerObj.Set("detour", "direct")
	dnsServersArr = append(dnsServersArr, dnsServerObj)
	dnsObj.Set("servers", dnsServersArr)
	config.Set("dns", dnsObj)
	// add inbound
	var inboundsArr serial.OrderedArray
	var socksObj serial.OrderedMap
	socksObj.Set("tag", "socks-in")
	socksObj.Set("listen", "::")
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
	var outbound2 serial.OrderedMap
	outbound2.Set("tag", "direct")
	outbound2.Set("type", "direct")
	outboundsArr = append(outboundsArr, outbound, outbound2)
	config.Set("outbounds", outboundsArr)
	// save test config
	marshal, err := json.Marshal(config)
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
	_ = os.Remove(configPath)
}
