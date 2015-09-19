package main

import (
	"flag"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/web"
	"log"
)

func main() {
	confiigLocation := flag.String("c", "", "set the json config file location")
	flag.Parse()
	log.Println(*confiigLocation)
	err := conf.LoadConfig(*confiigLocation)
	if err == nil {
		web.Start()
	} else {
		log.Fatalln(err)
	}
}
