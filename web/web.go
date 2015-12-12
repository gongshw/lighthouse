package web

import (
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/proxy"
	"github.com/gongshw/lighthouse/ssl"
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
	//http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	http.HandleFunc("/", staticHandler)
	var serverPortStr string
	if conf.CONFIG.DisableSSL {
		if conf.CONFIG.ServerPort == 0 {
			serverPortStr = strconv.Itoa(DEFAULT_HTTP_PORT)
		} else {
			serverPortStr = strconv.Itoa(conf.CONFIG.ServerPort)
		}
		log.Printf("http server listened at 0.0.0.0:%s", serverPortStr)
		log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
	} else {
		if conf.CONFIG.ServerPort == 0 {
			serverPortStr = strconv.Itoa(DEFAULT_HTTPS_PORT)
		} else {
			serverPortStr = strconv.Itoa(conf.CONFIG.ServerPort)
		}
		log.Printf("https server listened at 0.0.0.0:%s", serverPortStr)
		var certFile, keyFile string
		if conf.CONFIG.SSLCertificationFile != "" && conf.CONFIG.SSLKeyFile != "" {
			certFile = conf.CONFIG.SSLCertificationFile
			keyFile = conf.CONFIG.SSLKeyFile
		} else {
			var err error
			certFile, keyFile, err = ssl.CreateTempCrtAndKey(getHost())
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+serverPortStr, certFile, keyFile, nil))
	}
}

func getHost() string {
	url := conf.CONFIG.ServerBaseUrl
	if url == "" {
		return "localhost"
	} else {
		_, host := proxy.ParseBaseUrl(url)
		return host
	}
}
