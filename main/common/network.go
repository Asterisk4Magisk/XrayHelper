package common

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	tagNetwork = "network"
	timeout    = 3000
	dns        = "223.5.5.5:53"
)

// getHttpClient get an http client with custom dns
func getHttpClient(dns string, timeout time.Duration) *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Resolver: &net.Resolver{
					PreferGo: false,
					Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
						d := net.Dialer{Timeout: timeout}
						return d.DialContext(ctx, "udp", dns)
					},
				},
			}
			return dialer.DialContext(ctx, network, addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{Transport: transport}
}

func LookupIP(host string) ([]string, error) {
	resolver := &net.Resolver{
		PreferGo: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: timeout}
			return d.DialContext(ctx, "udp", dns)
		},
	}
	addrs, err := resolver.LookupIPAddr(context.Background(), host)
	if err != nil {
		return nil, e.New("lookup ipaddr failed, ", err)
	}
	ips := make([]string, len(addrs))
	for i, ia := range addrs {
		ips[i] = ia.IP.String()
	}
	return ips, nil
}

// CheckLocalPort check whether the local port is listening
func CheckLocalPort(pid string, port string, timeout time.Duration) bool {
	var check = func(ipv6 bool) bool {
		knetPath := "/proc/" + pid + "/net/tcp"
		if ipv6 {
			knetPath = "/proc/" + pid + "/net/tcp6"
		}
		i, _ := strconv.Atoi(port)
		// thx @young-zy, proc port always 4 characters hex width
		hex := fmt.Sprintf(":%04X ", i)
		if knet, err := os.ReadFile(knetPath); err == nil {
			return strings.Contains(string(knet), hex)
		}
		return false
	}
	start := time.Now()
	for time.Since(start) < timeout {
		if check(true) || check(false) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func IsIPv6(cidr string) bool {
	ip, _, _ := net.ParseCIDR(cidr)
	if ip != nil && ip.To4() == nil {
		return true
	}
	return false
}

func CheckLocalDevice(dev string, timeout time.Duration) bool {
	start := time.Now()
	for time.Since(start) < timeout {
		devices, err := net.Interfaces()
		if err == nil {
			for _, device := range devices {
				if dev == device.Name {
					return true
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// DownloadFile download file from url, and save to filepath
func DownloadFile(filepath string, url string) error {
	// get file from url
	client := getHttpClient(dns, timeout*time.Millisecond)
	request, _ := http.NewRequest("GET", url, nil)
	if len(builds.Config.XrayHelper.UserAgent) > 0 {
		request.Header.Set("User-Agent", builds.Config.XrayHelper.UserAgent)
	}
	response, err := client.Do(request)
	if err != nil {
		return e.New("cannot get file "+url+", ", err).WithPrefix(tagNetwork)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return e.New("bad http status "+response.Status+", ", err).WithPrefix(tagNetwork)
	}
	// open saveFile
	saveFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
	if err != nil {
		return e.New("cannot open file "+filepath+", ", err).WithPrefix(tagNetwork)
	}
	defer func(saveFile *os.File) {
		_ = saveFile.Close()
	}(saveFile)
	_, err = io.Copy(saveFile, response.Body)
	if err != nil {
		return e.New("save file "+filepath+" failed, ", err).WithPrefix(tagNetwork)
	}
	return nil
}

// GetRawData get raw data from a url
func GetRawData(url string) ([]byte, error) {
	client := getHttpClient(dns, timeout*time.Millisecond)
	request, _ := http.NewRequest("GET", url, nil)
	if len(builds.Config.XrayHelper.UserAgent) > 0 {
		request.Header.Set("User-Agent", builds.Config.XrayHelper.UserAgent)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, e.New("cannot get url "+url+", ", err).WithPrefix(tagNetwork)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, e.New("bad http status "+response.Status+", ", err).WithPrefix(tagNetwork)
	}
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.New("read data failed, ", err).WithPrefix(tagNetwork)
	}
	return raw, nil
}
