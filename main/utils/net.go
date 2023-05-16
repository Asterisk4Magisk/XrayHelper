package utils

import (
	"XrayHelper/main/errors"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// CheckPort check whether the port is listening
func CheckPort(protocol string, host string, port string) bool {
	addr := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout(protocol, addr, 1*time.Second)
	if err != nil {
		return false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	return true
}

// GetIPv6Addr get external ipv6 address, which should bypass
func GetIPv6Addr() ([]string, error) {
	var ipv6Addrs []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("cannot get ip address from local interface, ", err).WithPrefix("net")
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			if ipnet.IP.To4() == nil {
				ipv6Addrs = append(ipv6Addrs, ipnet.IP.String())
			}
		}
	}
	return ipv6Addrs, nil
}

// DownloadFile download file from url, and save to filepath
func DownloadFile(filepath string, url string) error {
	saveFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("cannot open file "+filepath+", ", err).WithPrefix("net")
	}
	defer func(saveFile *os.File) {
		_ = saveFile.Close()
	}(saveFile)
	response, err := http.Get(url)
	if err != nil {
		return errors.New("cannot get file "+url+", ", err).WithPrefix("net")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return errors.New("bad http status "+response.Status+", ", err).WithPrefix("net")
	}
	_, err = io.Copy(saveFile, response.Body)
	if err != nil {
		return errors.New("save file "+filepath+" failed, ", err).WithPrefix("net")
	}
	return nil
}
