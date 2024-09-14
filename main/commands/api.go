package commands

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/switches"
	"encoding/json"
	"fmt"
	"strconv"
)

const tagApi = "api"

type API struct {
	Operation string
	Object    string
	Addon     []string
}

type ApiCommand struct{}

func load() error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if err := builds.LoadPackage(); err != nil {
		return err
	}
	return nil
}

func (this *ApiCommand) Execute(args []string) error {
	if len(args) == 0 {
		fmt.Println(builds.Version())
		return nil
	} else if len(args) < 2 {
		return nil
	}
	api := API{Operation: args[0], Object: args[1], Addon: args[2:]}
	response, err := json.Marshal(parse(&api))
	if err == nil {
		fmt.Println(string(response))
	} else {
		err = e.New("api internal error, ", err).WithPrefix(tagApi).WithPathObj(*this)
	}
	return err
}

func parse(api *API) (response *serial.OrderedMap) {
	response = new(serial.OrderedMap)
	if err := load(); err != nil {
		return
	}
	switch api.Operation {
	case "get":
		switch api.Object {
		case "status":
			getStatus(response)
		case "switch":
			getSwitch(api, response)
		}
	case "set":
		switch api.Object {
		case "switch":
			setSwitch(api, response)
		}
	case "misc":
		switch api.Object {
		case "realping":
			realPing(api, response)
		}
	}
	return
}

func getStatus(response *serial.OrderedMap) {
	response.Set("api", builds.Version())
	response.Set("coreType", builds.Config.XrayHelper.CoreType)
	response.Set("pid", getServicePid())
	response.Set("method", builds.Config.Proxy.Method)
	response.Set("dataDir", builds.Config.XrayHelper.DataDir)
}

func getSwitch(api *API, response *serial.OrderedMap) {
	custom := false
	if len(api.Addon) > 0 && api.Addon[0] == "custom" {
		custom = true
	}
	if s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType); err == nil {
		var result serial.OrderedArray
		for _, url := range s.Get(custom) {
			result = append(result, url)
		}
		response.Set("result", result)
	}
}

func setSwitch(api *API, response *serial.OrderedMap) {
	response.Set("ok", false)
	custom := false
	index := 0
	if len(api.Addon) == 2 && api.Addon[0] == "custom" {
		custom = true
		index, _ = strconv.Atoi(api.Addon[1])
	} else if len(api.Addon) == 1 {
		index, _ = strconv.Atoi(api.Addon[0])
	} else {
		return
	}
	if s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType); err == nil {
		if err := s.Set(custom, index); err == nil {
			// if core is running, restart it
			if len(getServicePid()) > 0 {
				if err := restartService(); err == nil {
					response.Set("ok", true)
				}
			} else {
				response.Set("ok", true)
			}
		}
	}
}

func realPing(api *API, response *serial.OrderedMap) {
	response.Set("result", -1)
	custom := false
	index := 0
	if len(api.Addon) == 2 && api.Addon[0] == "custom" {
		custom = true
		index, _ = strconv.Atoi(api.Addon[1])
	} else if len(api.Addon) == 1 {
		index, _ = strconv.Atoi(api.Addon[0])
	} else {
		return
	}
	if s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType); err == nil {
		if target := s.Choose(custom, index); target != nil {
			if url, ok := target.(shareurls.ShareUrl); ok {
				response.Set("result", shareurls.RealPing(builds.Config.XrayHelper.CoreType, url))
			}
		}
	}
}
