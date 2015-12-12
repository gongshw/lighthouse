package web

import (
	"errors"
	"fmt"
	"github.com/gongshw/lighthouse/bindata"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/proxy"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	resp, conErr := proxy.ProxyRequest(r)
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
	if err:= proxy.ProxyResponse(w, resp, url); err!=nil{
		ShowError(w, "read content error", url)
	}
}

func detecthander(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.RequestURI()
	url, pathErr := proxyUrl(uri)
	if pathErr != nil {
		ShowError(w, "path error", uri)
	} else if !UrlNeedProxy(url) {
		ShowError(w, "url  blocked by admin", url)
	} else {
		http.Redirect(w, r, proxy.GetProxiedUrl(url, ""), http.StatusTemporaryRedirect)
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

func ShowError(w http.ResponseWriter, msg string, url string) {
	t := template.New("error_template")
	if t, err := t.Parse(string(bindata.MustAsset("error.html"))); err == nil {
		t.Execute(w, map[string]string{"msg": msg, "url": url})
	} else {
		log.Println("can't parse error.html")
		fmt.Fprintf(w, msg+": %s", url)
	}
}
