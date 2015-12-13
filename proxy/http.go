package proxy

import (
	"errors"
	"github.com/gongshw/lighthouse/conf"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func ProxyRequest(r *http.Request) (*http.Response, error) {
	// TODO proxy COOKIE
	url, _ := proxyUrl(r.URL.RequestURI())
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return nil, err
	}
	for k, vs := range r.Header {
		if isReqHeaderIgnore(k) {
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

func ProxyResponse(w http.ResponseWriter, resp *http.Response, url string) error {
	proxyHeader(resp.Header, w.Header(), url);
	w.WriteHeader(resp.StatusCode)
	var body []byte
	var readErr error
	if headerIs(resp.Header, "Content-Type", "text/html") {
		body, readErr = ParseHtml(resp.Body, url)
	} else if headerIs(resp.Header, "Content-Type", "text/css") {
		body, _ = ioutil.ReadAll(resp.Body)
		body = []byte(ParseCss(string(body[:]), url))
	} else {
		body, readErr = ioutil.ReadAll(resp.Body)
	}
	if readErr == nil {
		w.Write(body)
		return nil
	} else {
		log.Println(readErr)
		return errors.New("read content error")
	}
}

func proxyHeader(raw http.Header, proxied http.Header, url string){
	proxied.Add("Proxy-By", "gongshw/lighthouse")
	for key, valueArray := range raw {
		if isRespHeaderIgnore(key) {
			//ignore
		} else if key == "Location" {
			proxied.Set(key, GetProxiedUrl(raw.Get(key), url))
		} else {
			for _, value := range valueArray {
				proxied.Add(key, value)
			}
		}
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

func isReqHeaderIgnore(headName string) bool {
	switch headName{
		case
			"Cookie", // not support cookie
			"Accept-Encoding": // not support gzip encoding
			return true;
	}
	return false;
}

func isRespHeaderIgnore(headName string) bool {
	switch headName{
		case
			"Set-Cookie",
			"Content-Length",
			"Content-Security-Policy":
			return true;
	}
	return false;
}

func headerIs(headerMap map[string][]string, headerKey string, headerValue string) bool {
	header, exist := headerMap[headerKey]
	return exist && len(header) == 1 && strings.HasPrefix(header[0], headerValue)
}
