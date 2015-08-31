package hook

import (
	"github.com/gongshw/lighthouse/conf"
	"log"
	urllib "net/url"
	"regexp"
	"strings"
)

func ParseBaseUrl(url string) (string, string) {
	u, err := urllib.Parse(url)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return u.Scheme, u.Host
}

func ParseUrl(url string) (string, string) {
	u, err := urllib.Parse(url)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return u.Scheme, u.Host + u.RequestURI()
}

func GetProxiedUrl(url string) string {
	var serverBase string = conf.CONFIG.ServerBaseUrl
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		protocal, uri := ParseUrl(url)
		return serverBase + "/proxy/" + protocal + "/" + uri
	} else if strings.HasPrefix(url, "/") {
		return "/proxy" + url
	} else {
		return url
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
