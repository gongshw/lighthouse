package main

import (
	"flag"
	"fmt"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/web"
	"log"
	"os"
)

var (
	Version   string
	BuildTime string
)
var configLocation = flag.String("config", "", "Set the json config file location. If this flag is not set, lighthouse will lookup for a lighthouse.json file from the current direction and the $HOME of current user. If no config file found, lighthouse will use a default.")

var useDefaltConfig = flag.Bool("default", false, "Ignore all config location setting. Use the default config.")

var showVersion = flag.Bool("version", false, "Pring version info and exit with status code 0.")

func main() {
	flag.Parse()
	var err error
	if *showVersion {
		printVersionAndExit()
	}
	if !*useDefaltConfig {
		err = conf.LoadConfig(*configLocation)
	} else {
		log.Println("use default config")
	}
	if err == nil {
		web.Start()
	} else {
		log.Fatalln(err)
	}
}

func printVersionAndExit() {
	fmt.Println("Version: " + Version)
	fmt.Println("Build At: " + BuildTime)
	os.Exit(0)
}
