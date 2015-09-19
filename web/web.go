package web

import (
	"github.com/gongshw/lighthouse/conf"
	"log"
	"net/http"
	"strconv"
)

const (
	DEFAULT_HTTP_PORT  = 8080
	DEFAULT_HTTPS_PORT = 8443
)

func Start() {
	if initFilterErr := InitFilter(); initFilterErr != nil {
		log.Fatal(initFilterErr)
	}
	http.HandleFunc("/proxy/", proxyHandler)
	http.HandleFunc("/detect/", detecthander)
	http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	if conf.CONFIG.DisableSSL {
		var serverPortStr string
		if conf.CONFIG.ServerPort == 0 {
			serverPortStr = strconv.Itoa(DEFAULT_HTTP_PORT)
		} else {
			serverPortStr = strconv.Itoa(conf.CONFIG.ServerPort)
		}
		log.Printf("http server listened at 0.0.0.0:%s", serverPortStr)
		log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
	} else {
		log.Fatalln("https not supported")
	}
}
