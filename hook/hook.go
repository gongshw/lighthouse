package hook

import (
	"log"
	urllib "net/url"
	"regexp"
	"strings"
)

func GetBaseFromUrl(url string) string {
	u, err := urllib.Parse(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func GetResouceDir(url string) string {
	if strings.HasSuffix(url, "/") {
		return url
	} else {
		fileNameRegex := regexp.MustCompile("\\/[^\\/]+$")
		return fileNameRegex.ReplaceAllString(url, "/")
	}
}
