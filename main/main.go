package main

import (
	"flag"
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	configPath  string
	startXray   bool
	showVersion bool
)

func main() {
	fmt.Println(VersionStatement())
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&configPath, "c", "./config.yaml", "config file path")
	flag.Parse()
	if startXray {
		fmt.Println(Version())
		return
	}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := defaults.Set(&config); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config.XrayHelper)
	fmt.Printf("%+v\n", config.Proxy)
}
