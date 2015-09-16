package hook

import (
	"github.com/gongshw/lighthouse/conf"
	"log"
	urllib "net/url"
	"regexp"
	"strings"
)

func ParseBaseUrl(url string) (string, string) {
	u, err := urllib.ParseRequestURI(url)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return u.Scheme, u.Host
}

func ParseUrl(url string) (string, string) {
	u, err := urllib.ParseRequestURI(url)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return u.Scheme, u.Host + u.RequestURI()
}

func GetProxiedUrl(url string, base string) string {
	if !urlNeedProxy(url) {
		return url
	}
	var serverBase string = conf.CONFIG.ServerBaseUrl
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		protocal, uri := ParseUrl(url)
		return serverBase + "/proxy/" + protocal + "/" + uri
	} else {
		protocal, host := ParseBaseUrl(base)
		if strings.HasPrefix(url, "//") {
			return GetProxiedUrl(protocal+":"+url, base)
		} else if strings.HasPrefix(url, "/") {
			return serverBase + "/proxy/" + protocal + "/" + host + url
		} else {
			// a relative url
			return url
		}
	}
	return url
}

func urlNeedProxy(url string) bool {
	if strings.HasPrefix(url, "/proxy/http") {
		return false
	} else {
		return true
	}
}

func GetResouceDir(url string) string {
	if strings.HasSuffix(url, "/") {
		return url
	} else {
		fileNameRegex := regexp.MustCompile("\\/[^\\/]+$")
		return fileNameRegex.ReplaceAllString(url, "/")
	}
}
