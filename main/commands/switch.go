package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/shareurls"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tailscale/hujson"
	"os"
	"path"
	"strings"
)

var shareUrls []shareurls.ShareUrl

type SwitchCommand struct{}

func (this *SwitchCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) != 0 {
		return errors.New("too many arguments").WithPrefix("switch").WithPathObj(*this)
	}
	if err := loadShareUrl(); err != nil {
		return err
	}
	index := 0
	successFlag := false
	printProxyNode()
	fmt.Print("Please choose a node: ")
	_, err := fmt.Scanln(&index)
	if err != nil {
		return errors.New("invalid input, ", err).WithPrefix("switch").WithPathObj(*this)
	}
	if index < 0 || index >= len(shareUrls) {
		return errors.New("invalid node number").WithPrefix("switch").WithPathObj(*this)
	}
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return errors.New("open core config file failed, ", err).WithPrefix("switch").WithPathObj(*this)
	} else {
		if confInfo.IsDir() {
			confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return errors.New("open config dir failed, ", err).WithPrefix("switch").WithPathObj(*this)
			}
			for _, conf := range confDir {
				if !conf.IsDir() {
					confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()))
					if err != nil {
						return errors.New("read config file failed, ", err).WithPrefix("switch").WithPathObj(*this)
					}
					newConfByte, err := replaceProxyNode(confByte, index)
					if err != nil {
						log.HandleDebug(err)
						continue
					}
					if err := os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), newConfByte, 0644); err != nil {
						return errors.New("write new config failed, ", err).WithPrefix("service").WithPathObj(*this)
					}
					successFlag = true
				}
			}
		} else {
			confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return errors.New("read config file failed, ", err).WithPrefix("switch").WithPathObj(*this)
			}
			newConfByte, err := replaceProxyNode(confByte, index)
			if err != nil {
				return err
			}
			if err := os.WriteFile(builds.Config.XrayHelper.CoreConfig, newConfByte, 0644); err != nil {
				return errors.New("write new config failed, ", err).WithPrefix("service").WithPathObj(*this)
			}
			successFlag = true
		}
	}
	if successFlag {
		log.HandleInfo("switch: switch proxy node success")
	} else {
		return errors.New("switch proxy node failed").WithPrefix("service").WithPathObj(*this)
	}
	// if core is running, restart it
	if len(getServicePid()) > 0 {
		log.HandleInfo("switch: detect core is running, restart it")
		stopService()
		if err := startService(); err != nil {
			log.HandleError("restart service failed, " + err.Error())
		}
	}
	return nil
}

func loadShareUrl() error {
	subFile, err := os.Open(path.Join(builds.Config.XrayHelper.DataDir, "sub.txt"))
	if err != nil {
		return errors.New("open subscribe file failed, ", err).WithPrefix("switch")
	}
	defer func(subFile *os.File) {
		_ = subFile.Close()
	}(subFile)
	subScanner := bufio.NewScanner(subFile)
	subScanner.Split(bufio.ScanLines)
	for subScanner.Scan() {
		url := strings.TrimSpace(subScanner.Text())
		if len(url) > 0 {
			shareUrl, err := shareurls.NewShareUrl(url)
			if err != nil {
				log.HandleInfo("switch: " + err.Error() + ", drop it")
				continue
			}
			shareUrls = append(shareUrls, shareUrl)
		}
	}
	if len(shareUrls) == 0 {
		return errors.New("no valid nodes").WithPrefix("switch")
	}
	return nil
}

func printProxyNode() {
	for index, shareUrl := range shareUrls {
		fmt.Printf("[%d] %s\n", index, shareUrl.GetNodeInfo())
	}
}

func replaceProxyNode(conf []byte, index int) (replacedConf []byte, err error) {
	// standardize origin json (remove comment)
	standardize, err := hujson.Standardize(conf)
	if err != nil {
		return nil, errors.New("standardize config json failed, ", err).WithPrefix("switch")
	}
	// unmarshal
	var jsonValue interface{}
	err = json.Unmarshal(standardize, &jsonValue)
	if err != nil {
		return nil, errors.New("unmarshal config json failed, ", err).WithPrefix("switch")
	}
	// assert json to map
	jsonMap, ok := jsonValue.(map[string]interface{})
	if !ok {
		return nil, errors.New("assert config json to map failed").WithPrefix("switch")
	}
	outbounds, ok := jsonMap["outbounds"]
	if !ok {
		return nil, errors.New("cannot find outbounds ").WithPrefix("switch")
	}
	// assert outbounds
	outboundsMap, ok := outbounds.([]interface{})
	if !ok {
		return nil, errors.New("assert outbounds to []interface failed, ").WithPrefix("switch")
	}
	for i, outbound := range outboundsMap {
		outboundMap, ok := outbound.(map[string]interface{})
		if !ok {
			continue
		}
		tag, ok := outboundMap["tag"].(string)
		if !ok {
			continue
		}
		if tag == "proxy" {
			// replace
			outbound, err = shareUrls[index].ToOutoundWithTag(builds.Config.XrayHelper.CoreType, builds.Config.XrayHelper.ProxyTag)
			if err != nil {
				return nil, err
			}
			outboundsMap[i] = outbound
			jsonMap["outbounds"] = outboundsMap
			// marshal
			marshal, err := json.MarshalIndent(jsonMap, "", "    ")
			if err != nil {
				return nil, errors.New("marshal config json failed, ", err).WithPrefix("switch")
			}
			return marshal, nil
		}
	}
	return nil, errors.New("not found tag, " + builds.Config.XrayHelper.ProxyTag).WithPrefix("switch")
}
