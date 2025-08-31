package shareurls

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls/addon"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

const (
	tagSpeedtest = "speedtest"
	testUrl      = "https://www.google.com/generate_204"
)

type Result struct {
	Index string
	Url   ShareUrl
	Port  int
	Value int
}

func RealPing(coreType string, results []*Result) {
	configPath := path.Join(builds.Config.XrayHelper.RunDir, "test.json")
	// start test service
	service, err := startTestService(coreType, configPath, results)
	if err != nil {
		log.HandleDebug(err)
		return
	}
	defer stopTestService(service, configPath)
	// check service port
	for _, result := range results {
		if common.CheckLocalPort(strconv.Itoa(service.Pid()), strconv.Itoa(result.Port), 2*time.Second) {
			continue
		}
		log.HandleDebug("service not listen for RealPing")
		return
	}
	var wg sync.WaitGroup
	for _, result := range results {
		wg.Add(1)
		go func(result *Result) {
			defer wg.Done()
			// set socks5 proxy
			dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:"+strconv.Itoa(result.Port), nil, proxy.Direct)
			if err != nil {
				log.HandleDebug("set socks5 proxy: " + err.Error())
				return
			}
			start := time.Now()
			for {
				result.Value = startTest(dialer)
				if time.Since(start) > 4*time.Second || result.Value > -1 {
					break
				}
			}
			return
		}(result)
	}
	wg.Wait()
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
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
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

func startTestService(coreType string, configPath string, results []*Result) (common.External, error) {
	var service common.External
	switch coreType {
	case "xray":
		if err := genXrayTestConfig(configPath, results); err != nil {
			return nil, err
		}
		service = common.NewExternal(0, nil, nil, builds.Config.XrayHelper.CorePath, "run", "-c", configPath)
	case "sing-box":
		if err := genSingboxTestConfig(configPath, results); err != nil {
			return nil, err
		}
		service = common.NewExternal(0, nil, nil, builds.Config.XrayHelper.CorePath, "run", "-c", configPath, "--disable-color")
	default:
		return nil, e.New("not a supported coreType " + coreType).WithPrefix(tagSpeedtest)
	}
	service.SetUidGid("0", common.CoreGid)
	service.Start()
	if service.Err() != nil {
		return nil, e.New("start test service failed, ", service.Err()).WithPrefix(tagSpeedtest)
	}
	return service, nil
}

func genXrayTestConfig(configPath string, results []*Result) error {
	var nodeInfo []*addon.NodeInfo
	for _, result := range results {
		nodeInfo = append(nodeInfo, result.Url.GetNodeInfo())
	}
	var config serial.OrderedMap
	// add dns
	var dnsObj serial.OrderedMap
	var dnsHostsObj serial.OrderedMap
	for _, info := range nodeInfo {
		ip, err := common.LookupIP(info.Host)
		if err != nil {
			continue
		}
		dnsHostsObj.Set(info.Host, ip)
	}
	dnsObj.Set("hosts", dnsHostsObj)
	var dnsServersArr serial.OrderedArray
	dnsServersArr = append(dnsServersArr, "223.5.5.5")
	dnsObj.Set("servers", dnsServersArr)
	config.Set("dns", dnsObj)
	// add inbounds
	var inboundsArr serial.OrderedArray
	for _, result := range results {
		tag := "in-" + strconv.Itoa(result.Port)
		var socksObj serial.OrderedMap
		socksObj.Set("tag", tag)
		socksObj.Set("port", result.Port)
		socksObj.Set("protocol", "socks")

		var sniffingObj serial.OrderedMap
		sniffingObj.Set("enabled", true)
		var destOverrideArr serial.OrderedArray
		destOverrideArr = append(destOverrideArr, "http", "tls", "quic")
		sniffingObj.Set("destOverride", destOverrideArr)

		socksObj.Set("sniffing", sniffingObj)
		inboundsArr = append(inboundsArr, socksObj)
	}
	config.Set("inbounds", inboundsArr)
	// add outbounds
	var outboundsArr serial.OrderedArray
	for i, result := range results {
		tag := "out-" + strconv.Itoa(i)
		outbound, err := result.Url.ToOutboundWithTag("xray", tag)
		if err != nil {
			return err
		}
		outboundsArr = append(outboundsArr, outbound)
	}
	config.Set("outbounds", outboundsArr)
	// add routing
	var routing serial.OrderedMap
	var rulesArr serial.OrderedArray
	for i, result := range results {
		inTag := "in-" + strconv.Itoa(result.Port)
		outTag := "out-" + strconv.Itoa(i)
		var rule serial.OrderedMap
		rule.Set("type", "field")
		var inboundTag serial.OrderedArray
		inboundTag = append(inboundTag, inTag)
		rule.Set("inboundTag", inboundTag)
		rule.Set("outboundTag", outTag)
		rulesArr = append(rulesArr, rule)
	}
	routing.Set("rules", rulesArr)
	config.Set("routing", routing)
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

func genSingboxTestConfig(configPath string, results []*Result) error {
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
	for _, result := range results {
		tag := "in-" + strconv.Itoa(result.Port)
		var socksObj serial.OrderedMap
		socksObj.Set("tag", tag)
		socksObj.Set("listen", "::")
		socksObj.Set("listen_port", result.Port)
		socksObj.Set("type", "socks")
		socksObj.Set("sniff", true)
		socksObj.Set("sniff_override_destination", true)
		inboundsArr = append(inboundsArr, socksObj)
	}
	config.Set("inbounds", inboundsArr)
	// add outbounds
	var outboundsArr serial.OrderedArray
	var direct serial.OrderedMap
	direct.Set("tag", "direct")
	direct.Set("type", "direct")
	outboundsArr = append(outboundsArr, direct)
	for i, result := range results {
		tag := "out-" + strconv.Itoa(i)
		outbound, err := result.Url.ToOutboundWithTag("sing-box", tag)
		if err != nil {
			return err
		}
		outboundsArr = append(outboundsArr, outbound)
	}
	config.Set("outbounds", outboundsArr)
	// add route
	var route serial.OrderedMap
	route.Set("final", "direct")
	var rulesArr serial.OrderedArray
	for i, result := range results {
		inTag := "in-" + strconv.Itoa(result.Port)
		outTag := "out-" + strconv.Itoa(i)
		var rule serial.OrderedMap
		var inbound serial.OrderedArray
		inbound = append(inbound, inTag)
		rule.Set("inbound", inbound)
		rule.Set("outbound", outTag)
		rulesArr = append(rulesArr, rule)
	}
	route.Set("rules", rulesArr)
	config.Set("route", route)
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
