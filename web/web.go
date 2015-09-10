package web

import (
	"errors"
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

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.RequestURI()
	if pathUrlHostOnly(uri) {
		http.Redirect(w, r, uri+"/", http.StatusTemporaryRedirect)
		return
	}
	url, pathErr := proxyUrl(uri)
	if pathErr != nil {
		log.Println(pathErr)
		ShowError(w, "path error", uri)
		return
	}
	if !UrlNeedProxy(url) {
		log.Printf("no need to proxy: %s", url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	resp, conErr := proxyRequest(r)
	if conErr != nil {
		log.Printf("connection error:[%s,%s]", conErr, url)
		ShowError(w, "connection error", url)
		return
	}
	defer resp.Body.Close()
	if size := resp.ContentLength; size > conf.CONFIG.ContentLengthLimit {
		ShowError(w, "responese too large", url)
		return
	}
	log.Printf("%s %s %s", r.Method, url, resp.Status)
	proxyResponse(w, resp, url)
}

func detecthander(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.RequestURI()
	url, pathErr := proxyUrl(uri)
	if pathErr != nil {
		ShowError(w, "path error", uri)
	} else if !UrlNeedProxy(url) {
		ShowError(w, "url  blocked by admin", url)
	} else {
		http.Redirect(w, r, hook.GetProxiedUrl(url, ""), http.StatusTemporaryRedirect)
	}
}

func proxyUrl(path string) (string, error) {
	token := strings.SplitN(path, "/", 5)
	if len(token) == 5 {
		return token[2] + "://" + token[3] + "/" + token[4], nil
	} else if len(token) == 4 {
		return token[2] + "://" + token[3] + "/", nil
	}
	return "", errors.New("illegal path: " + path)
}

func pathUrlHostOnly(uri string) bool {
	return strings.Count(uri, "/") < 4
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

func proxyResponse(w http.ResponseWriter, resp *http.Response, url string) {
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

	var body []byte
	var readError error
	if headerIs(resp.Header, "Content-Type", "text/html") {
		body, readErr := hook.ParseHtml(resp.Body, url)
	} else {
		body, readErr := ioutil.ReadAll(resp.Body)
	}
	if readErr == nil {
		w.Write(body)
	} else {
		log.Println(readErr)
		ShowError(w, "read content error", url)
	}
}

func headerIs(headerMap map[string][]string, headerKey string, headerValue string) bool {
	header, exist := headerMap[headerKey]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], headerValue)
}

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
