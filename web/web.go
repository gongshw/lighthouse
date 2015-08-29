package web

import (
	"errors"
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

const _5MB = 5 * 1024 * 1024

var ERROR_RESPONSE_TOO_LARGE = errors.New("responese too large")

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.RawQuery
	resp, conErr := proxyRequest(r)
	if conErr != nil {
		log.Println(conErr)
		fmt.Fprintf(w, "connection error: %s", url)
		return
	}
	defer resp.Body.Close()
	if size := resp.ContentLength; size > _5MB {
		log.Fatal(ERROR_RESPONSE_TOO_LARGE)
		fmt.Fprintf(w, "responese too large: %s", url)
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
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		fmt.Fprintf(w, "ead content error: %s", url)
		return
	}
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
	log.Println("proxy " + url)
	resp, conErr := http.Get(url)
	return resp, conErr
}

func headerIs(headerMap map[string][]string, headerKey string, headerValue string) bool {
	header, exist := headerMap[headerKey]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], headerValue)
}

func Start() {
	http.HandleFunc("/proxy", proxyHandler)
	http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	serverPortStr := strconv.Itoa(conf.CONFIG.ServerPort)
	log.Println("server listened at 0.0.0.0:" + serverPortStr)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
}
