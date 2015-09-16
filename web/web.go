package web

import (
	"github.com/gongshw/lighthouse/conf"
	"log"
	"net/http"
	"strconv"
)

func Start() {
	if initFilterErr := InitFilter(); initFilterErr != nil {
		log.Fatal(initFilterErr)
	}
	http.HandleFunc("/proxy/", proxyHandler)
	http.HandleFunc("/detect/", detecthander)
	http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	serverPortStr := strconv.Itoa(conf.CONFIG.ServerPort)
	log.Printf("server listened at 0.0.0.0:%s", serverPortStr)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
}
