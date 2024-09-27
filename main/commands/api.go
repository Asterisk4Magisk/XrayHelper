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

func (this *ApiCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
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
	get := func(custom bool) serial.OrderedArray {
		var result serial.OrderedArray
		if s, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType); err == nil {
			defer s.Clear()
			for _, url := range s.Get(custom) {
				result = append(result, url)
			}
		}
		return result
	}
	if len(api.Addon) > 0 {
		if api.Addon[0] == "all" {
			response.Set("result", get(false))
			response.Set("custom", get(true))
		} else if api.Addon[0] == "custom" {
			response.Set("result", get(true))
		}
	} else {
		response.Set("result", get(false))
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
	var responseArr serial.OrderedArray
	response.Set("result", responseArr)
	if len(api.Addon) == 0 {
		return
	}
	start := func(index []string, custom bool) (arr serial.OrderedArray) {
		var (
			results []*shareurls.Result
			res     []*shareurls.Result
			port    = 65500
			i       = 0
		)
		if swh, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType); err == nil {
			for _, idx := range index {
				id, _ := strconv.Atoi(idx)
				if target := swh.Choose(custom, id); target != nil {
					if url, ok := target.(shareurls.ShareUrl); ok {
						if i > 50 {
							shareurls.RealPing(builds.Config.XrayHelper.CoreType, res)
							results = append(results, res...)
							res = make([]*shareurls.Result, 0)
							port = 65500
							i = 0
						}
						res = append(res, &shareurls.Result{Index: idx, Url: url, Port: port, Value: -1})
						port -= 1
						i++
					}
				}
			}
		}
		shareurls.RealPing(builds.Config.XrayHelper.CoreType, res)
		results = append(results, res...)
		for _, result := range results {
			var ret serial.OrderedMap
			ret.Set("index", result.Index)
			ret.Set("realping", result.Value)
			arr = append(arr, ret)
		}
		return
	}
	if api.Addon[0] == "custom" {
		response.Set("result", start(api.Addon[1:], true))
	} else {
		response.Set("result", start(api.Addon, false))
	}
}
