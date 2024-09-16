package ray

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

const tagRayswitch = "rayswitch"

var shareUrls []shareurls.ShareUrl

type RaySwitch struct{}

func (this *RaySwitch) Execute(args []string) (bool, error) {
	if len(args) > 1 {
		return false, e.New("too many arguments").WithPrefix(tagRayswitch).WithPathObj(*this)
	}
	if len(args) == 1 && args[0] == "custom" {
		if err := loadShareUrl(true); err != nil {
			return false, err
		}
	} else {
		if err := loadShareUrl(false); err != nil {
			return false, err
		}
	}
	printProxyNode()
	fmt.Print("Please choose a node: ")
	index := 0
	if _, err := fmt.Scanln(&index); err != nil {
		return false, e.New("invalid input, ", err).WithPrefix(tagRayswitch).WithPathObj(*this)
	}
	if err := change(index); err != nil {
		return false, err
	}
	return true, nil
}

func (this *RaySwitch) Get(custom bool) serial.OrderedArray {
	var result serial.OrderedArray
	err := loadShareUrl(custom)
	if err == nil {
		for _, url := range shareUrls {
			result = append(result, url.GetNodeInfo())
		}
	}
	return result
}

func (this *RaySwitch) Set(custom bool, index int) error {
	err := loadShareUrl(custom)
	if err == nil {
		return change(index)
	}
	return err
}

func (this *RaySwitch) Choose(custom bool, index int) any {
	err := loadShareUrl(custom)
	if err == nil {
		if index >= 0 && index < len(shareUrls) {
			return shareUrls[index]
		}
	}
	return nil
}

func change(index int) error {
	if index < 0 || index >= len(shareUrls) {
		return e.New("invalid number").WithPrefix(tagRayswitch)
	}
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return e.New("open core config file failed, ", err).WithPrefix(tagRayswitch)
	} else {
		if confInfo.IsDir() {
			confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return e.New("open config dir failed, ", err).WithPrefix(tagRayswitch)
			}
			hostFlag := false
			replaceFlag := false
			if builds.Config.XrayHelper.CoreType != "xray" {
				hostFlag = true
			}
			for _, conf := range confDir {
				if !conf.IsDir() && strings.HasSuffix(conf.Name(), ".json") {
					confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()))
					if err != nil {
						log.HandleDebug("read config file failed, " + err.Error())
						continue
					}
					if !hostFlag {
						confByte, err = replaceXrayHost(confByte, index)
						if err != nil {
							log.HandleDebug(err)
						} else {
							err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), confByte, 0644)
							if err != nil {
								log.HandleDebug("write new config failed, " + err.Error())
							} else {
								hostFlag = true
							}
						}
					}
					if !replaceFlag {
						confByte, err = replaceProxyNode(confByte, index)
						if err != nil {
							log.HandleDebug(err)
						} else {
							err = os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), confByte, 0644)
							if err != nil {
								log.HandleDebug("write new config failed, " + err.Error())
							} else {
								replaceFlag = true
							}
						}
					}
					if hostFlag && replaceFlag {
						return nil
					}
				}
			}
		} else {
			confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return e.New("read config file failed, ", err).WithPrefix(tagRayswitch)
			}
			if builds.Config.XrayHelper.CoreType == "xray" {
				confByte, err = replaceXrayHost(confByte, index)
				if err != nil {
					return err
				}
			}
			confByte, err = replaceProxyNode(confByte, index)
			if err != nil {
				return err
			}
			if err := os.WriteFile(builds.Config.XrayHelper.CoreConfig, confByte, 0644); err != nil {
				return e.New("write new config failed, ", err).WithPrefix(tagRayswitch)
			}
			return nil
		}
	}
	return e.New("write new config failed").WithPrefix(tagRayswitch)
}

func loadShareUrl(custom bool) error {
	if len(shareUrls) > 0 {
		return nil
	}
	var nodeTxt string
	if custom {
		nodeTxt = path.Join(builds.Config.XrayHelper.DataDir, "custom.txt")
	} else {
		nodeTxt = path.Join(builds.Config.XrayHelper.DataDir, "sub.txt")
	}
	subFile, err := os.Open(nodeTxt)
	if err != nil {
		return e.New("open proxy node file failed, ", err).WithPrefix(tagRayswitch)
	}
	defer func(subFile *os.File) {
		_ = subFile.Close()
	}(subFile)
	subScanner := bufio.NewScanner(subFile)
	subScanner.Split(bufio.ScanLines)
	for subScanner.Scan() {
		url := strings.TrimSpace(subScanner.Text())
		if len(url) > 0 {
			shareUrl, err := shareurls.Parse(url)
			if err != nil {
				log.HandleDebug("switch: " + err.Error() + ", drop it")
				continue
			}
			shareUrls = append(shareUrls, shareUrl)
		}
	}
	if len(shareUrls) == 0 {
		return e.New("no valid nodes").WithPrefix(tagRayswitch)
	}
	return nil
}

func printProxyNode() {
	for index, shareUrl := range shareUrls {
		fmt.Printf(color.GreenString("[%d]")+" %s\n", index, shareUrl.GetNodeInfoStr())
	}
}

func replaceProxyNode(conf []byte, index int) ([]byte, error) {
	switch builds.Config.XrayHelper.CoreType {
	case "xray", "sing-box":
		// unmarshal
		var jsonMap serial.OrderedMap
		err := json.Unmarshal(conf, &jsonMap)
		if err != nil {
			return nil, e.New("unmarshal config json failed, ", err).WithPrefix(tagRayswitch)
		}
		if outbounds, ok := jsonMap.Get("outbounds"); ok {
			outboundArray := outbounds.Value.(serial.OrderedArray)
			for i, outbound := range outboundArray {
				outboundMap := outbound.(serial.OrderedMap)
				if tag, ok := outboundMap.Get("tag"); ok {
					if tag.Value == builds.Config.XrayHelper.ProxyTag {
						// replace
						outbound, err := shareUrls[index].ToOutboundWithTag(builds.Config.XrayHelper.CoreType, builds.Config.XrayHelper.ProxyTag)
						if err != nil {
							return nil, err
						}
						outboundArray[i] = outbound
						jsonMap.Set("outbounds", outboundArray)
						// marshal
						marshal, err := json.MarshalIndent(jsonMap, "", "    ")
						if err != nil {
							return nil, e.New("marshal config json failed, ", err).WithPrefix(tagRayswitch)
						}
						return marshal, nil
					}
				}
			}
			return nil, e.New("cannot found outbounds tag: " + builds.Config.XrayHelper.ProxyTag).WithPrefix(tagRayswitch)
		}
		return nil, e.New("cannot found outbounds from provided conf").WithPrefix(tagRayswitch)
	case "hysteria2":
		// unmarshal
		var yamlMap serial.OrderedMap
		err := yaml.Unmarshal(conf, &yamlMap)
		if err != nil {
			return nil, e.New("unmarshal config yaml failed, ", err).WithPrefix(tagRayswitch)
		}
		// get hysteria client config from shareUrl
		clientObject, err := shareUrls[index].ToOutboundWithTag(builds.Config.XrayHelper.CoreType, "")
		if err != nil {
			return nil, err
		}
		// replace
		if server, ok := clientObject.Get("server"); ok {
			yamlMap.SetValue(server)
		}
		if auth, ok := clientObject.Get("auth"); ok {
			yamlMap.SetValue(auth)
		}
		if obfs, ok := clientObject.Get("obfs"); ok {
			yamlMap.SetValue(obfs)
		}
		if tls, ok := clientObject.Get("tls"); ok {
			yamlMap.SetValue(tls)
		}
		// marshal
		marshal, err := yaml.Marshal(yamlMap)
		if err != nil {
			return nil, e.New("marshal config yaml failed, ", err).WithPrefix(tagRayswitch)
		}
		return marshal, nil
	}
	return nil, e.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix(tagRayswitch)
}

func replaceXrayHost(conf []byte, index int) ([]byte, error) {
	// unmarshal
	var jsonMap serial.OrderedMap
	err := json.Unmarshal(conf, &jsonMap)
	if err != nil {
		return nil, e.New("unmarshal config json failed, ", err).WithPrefix(tagRayswitch)
	}
	// asset dns
	if dns, ok := jsonMap.Get("dns"); ok {
		dnsMap := dns.Value.(serial.OrderedMap)
		// replace
		var hostsMap serial.OrderedMap
		nodeInfo := shareUrls[index].GetNodeInfo()
		result, err := common.LookupIP(nodeInfo.Host)
		if err != nil {
			return nil, err
		}
		hostsMap.Set(nodeInfo.Host, result)
		dnsMap.Set("hosts", hostsMap)
		jsonMap.Set("dns", dnsMap)
		// marshal
		marshal, err := json.MarshalIndent(jsonMap, "", "    ")
		if err != nil {
			return nil, e.New("marshal config json failed, ", err).WithPrefix(tagRayswitch)
		}
		return marshal, nil
	}
	return nil, e.New("cannot find dns from provided conf").WithPrefix(tagRayswitch)
}
