package web

import (
	"encoding/json"
	"fmt"
	"github.com/gongshw/lighthouse/hook"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	resp, conErr := http.Get(url)
	if conErr != nil {
		fmt.Fprintf(w, "连接异常: %s", url)
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readErr != nil {
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

func headerIs(headerMap map[string][]string, contentType string, typrValue string) bool {
	header, exist := headerMap[contentType]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], typrValue)
}

func Start() {
	configFile, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	decoder := json.NewDecoder(configFile)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/proxy", proxyHandler)
	http.Handle("/", http.FileServer(http.Dir(configuration.StaicFileDir)))
	http.ListenAndServe(":"+strconv.Itoa(configuration.ServerPort), nil)
}
