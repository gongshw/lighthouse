package hook

import (
	"regexp"
)

var (
	urlPattern = regexp.MustCompile("(url\\(\\s*[\"']?\\s*)([^\\)]+)(\\s*[\"']?\\s*\\))")
)

func ParseCss(css string, url string) string {
	css = urlPattern.ReplaceAllStringFunc(css, getProxyUrlInCssFunc(url))
	return css
}

type rpl func(string) string

func getProxyUrlInCssFunc(baseUrl string) rpl {
	return func(token string) string {
		if match := urlPattern.FindStringSubmatch(token); len(match) == 4 {
			return match[1] + GetProxiedUrl(match[2], baseUrl) + match[3]
		} else {
			return token
		}
	}
}
