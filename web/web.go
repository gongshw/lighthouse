package web

import (
	"errors"
	"fmt"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/hook"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const _5MB = 5 * 1024 * 1024

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url, pathErr := proxyUrl(r.URL.RequestURI())
	if pathErr != nil {
		log.Println(pathErr)
		fmt.Fprintf(w, "path error: %s", r.URL.RawPath)
		return
	}
	resp, conErr := proxyRequest(r)
	if conErr != nil {
		log.Println(conErr)
		fmt.Fprintf(w, "connection error: %s", url)
		return
	}
	defer resp.Body.Close()
	if size := resp.ContentLength; size > _5MB {
		fmt.Fprintf(w, "responese too large: %s", url)
		return
	}
	log.Println(r.Method + " " + url + " " + resp.Status)
	w.Header().Add("Proxy-By", "gongshw/lighthouse")
	for key, valueArray := range resp.Header {
		if key == "Content-Length" || key == "Set-Cookie" {
			//ignore
		} else if key == "Location" {
			w.Header().Set(key, hook.GetProxiedUrl(resp.Header.Get(key), url))
		} else {
			for _, value := range valueArray {
				w.Header().Add(key, value)
			}
		}
	}
	w.WriteHeader(resp.StatusCode)
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		fmt.Fprintf(w, "read content error: %s", url)
		return
	}
	if headerIs(resp.Header, "Content-Type", "text/html") {
		w.Write([]byte(hook.ParseHtml(string(body[:]), url)))
	} else {
		w.Write(body)
	}
}

func proxyUrl(path string) (string, error) {
	token := strings.SplitN(path, "/", 5)
	if len(token) < 4 {
		return "", errors.New("illegal path: " + path)
	} else if len(token) == 4 {
		return token[2] + "://" + token[3], nil
	} else if len(token) == 5 {
		return token[2] + "://" + token[3] + "/" + token[4], nil
	}
	return "", errors.New("illegal path: " + path)
}

func proxyRequest(r *http.Request) (*http.Response, error) {
	// TODO proxy COOKIE
	url, _ := proxyUrl(r.URL.RequestURI())
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	for k, vs := range r.Header {
		if k == "Cookie" || k == "Accept-Encoding" {
			//ignore
		} else {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: conf.CONFIG.ResponseTimeoutSecond * time.Second,
		}).Dial,
	}
	return tr.RoundTrip(req)
}

func headerIs(headerMap map[string][]string, headerKey string, headerValue string) bool {
	header, exist := headerMap[headerKey]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], headerValue)
}

func Start() {
	http.HandleFunc("/proxy/", proxyHandler)
	http.Handle("/", http.FileServer(http.Dir(conf.CONFIG.StaicFileDir)))
	serverPortStr := strconv.Itoa(conf.CONFIG.ServerPort)
	log.Println("server listened at 0.0.0.0:" + serverPortStr)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+serverPortStr, nil))
}
