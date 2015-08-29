package main

import (
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/web"
	"log"
)

func main() {
	err := conf.LoadConfig("")
	if err == nil {
		web.Start()
	} else {
		log.Fatalln(err)
	}
}
