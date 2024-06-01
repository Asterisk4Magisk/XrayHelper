package ray

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
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

func (this *RaySwitch) Get(custom bool) []string {
	var result []string
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
			for _, conf := range confDir {
				if !conf.IsDir() {
					confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()))
					if err != nil {
						return e.New("read config file failed, ", err).WithPrefix(tagRayswitch)
					}
					newConfByte, err := replaceProxyNode(confByte, index)
					if err != nil {
						log.HandleDebug(err)
						continue
					}
					if err := os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), newConfByte, 0644); err != nil {
						return e.New("write new config failed, ", err).WithPrefix(tagRayswitch)
					}
					return nil
				}
			}
		} else {
			confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return e.New("read config file failed, ", err).WithPrefix(tagRayswitch)
			}
			newConfByte, err := replaceProxyNode(confByte, index)
			if err != nil {
				return err
			}
			if err := os.WriteFile(builds.Config.XrayHelper.CoreConfig, newConfByte, 0644); err != nil {
				return e.New("write new config failed, ", err).WithPrefix(tagRayswitch)
			}
			return nil
		}
	}
	return e.New("write new config failed, ").WithPrefix(tagRayswitch)
}

func loadShareUrl(custom bool) error {
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
				log.HandleInfo("switch: " + err.Error() + ", drop it")
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
		fmt.Printf(color.GreenString("[%d]")+" %s\n", index, shareUrl.GetNodeInfo())
	}
}

func replaceProxyNode(conf []byte, index int) (replacedConf []byte, err error) {
	// unmarshal
	var jsonMap serial.OrderedMap
	err = json.Unmarshal(conf, &jsonMap)
	if err != nil {
		return nil, e.New("unmarshal config json failed, ", err).WithPrefix(tagRayswitch)
	}
	outbounds, ok := jsonMap.Get("outbounds")
	if !ok {
		return nil, e.New("cannot find outbounds").WithPrefix(tagRayswitch)
	}
	// assert outbounds
	outboundArray, ok := outbounds.Value.(serial.OrderedArray)
	if !ok {
		return nil, e.New("assert outbounds to serial.OrderedArray failed").WithPrefix(tagRayswitch)
	}
	for i, outbound := range outboundArray {
		outboundMap, ok := outbound.(serial.OrderedMap)
		if !ok {
			continue
		}
		tag, ok := outboundMap.Get("tag")
		if !ok {
			continue
		}
		if tag.Value == builds.Config.XrayHelper.ProxyTag {
			// replace
			outbound, err = shareUrls[index].ToOutboundWithTag(builds.Config.XrayHelper.CoreType, builds.Config.XrayHelper.ProxyTag)
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
	return nil, e.New("not found tag, " + builds.Config.XrayHelper.ProxyTag).WithPrefix(tagRayswitch)
}
