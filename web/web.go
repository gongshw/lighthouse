package web

import (
	"fmt"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/hook"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Configuration struct {
	StaicFileDir  string
	ServerBaseUrl string
	ServerPort    int
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RawQuery
	resp, conErr := proxyRequest(r)
	if conErr != nil {
		log.Fatal(conErr)
		fmt.Fprintf(w, "连接异常: %s", url)
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
		fmt.Fprintf(w, "网页读取异常: %s", url)
		return
	}
	w.Header().Add("Proxy-By", "gongshw/lighthouse")
	for key, valueArray := range resp.Header {
		if key == "Content-Length" {
			continue
		}
		for _, value := range valueArray {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if headerIs(resp.Header, "Content-Type", "text/html") {
		w.Write([]byte(hook.ParseHtml(string(body[:]), url)))
	} else if headerIs(resp.Header, "Content-Type", "text/css") {
		w.Write([]byte(hook.ParseCss(string(body[:]), url)))
	} else {
		w.Write(body)
	}
}

func proxyRequest(r *http.Request) (*http.Response, error) {
	// TODO proxy METHOD and HEADER
	url := r.URL.RawQuery
	log.Println("Proxy " + url)
	resp, conErr := http.Get(url)
	return resp, conErr
}

func headerIs(headerMap map[string][]string, contentType string, typrValue string) bool {
	header, exist := headerMap[contentType]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], typrValue)
}

func Start() {
	http.HandleFunc("/proxy", proxyHandler)
	http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	serverPortStr := strconv.Itoa(conf.CONFIG.ServerPort)
	log.Println("Server listened at 0.0.0.0:" + serverPortStr)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
}
