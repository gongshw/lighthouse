package hook

import (
	"regexp"
)

func ParseCss(css string, url string) string {
	cssUrlRegex := regexp.MustCompile("url\\s*\\([\"']([\\w\\.\\/\\?#]+)[\"']\\s*\\)")
	dirPrefix := "/proxy?" + GetResouceDir(url)
	return cssUrlRegex.ReplaceAllString(css, "url(\""+dirPrefix+"$1"+"\")")
}
