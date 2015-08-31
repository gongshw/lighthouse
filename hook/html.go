package hook

import (
	"github.com/gongshw/lighthouse/conf"
	"regexp"
)

func ParseHtml(html string, url string) string {
	i := 0
	var htmlBuf []byte
	var tokenBuf []byte
	htmlLength := len(html)
	for {
		if i >= htmlLength {
			// end of html, flush tokenBuf
			if len(tokenBuf) > 0 {
				flushToken(&htmlBuf, tokenBuf, url)
			}
			break
		}
		b := html[i]
		if len(tokenBuf) == 0 {
			// new token
			tokenBuf = append(tokenBuf, b)
			i++
		} else if b == '<' {
			// new tag token start
			flushToken(&htmlBuf, tokenBuf, url)
			tokenBuf = tokenBuf[:0]
		} else {
			tokenBuf = append(tokenBuf, b)
			i++
			if b == '>' {
				flushToken(&htmlBuf, tokenBuf, url)
				tokenBuf = tokenBuf[:0]
			}
		}
	}
	return string(htmlBuf)
}

func flushToken(htmlBuf *[]byte, tokenBuf []byte, url string) {
	var serverBase string = conf.CONFIG.ServerBaseUrl
	var JS_HOOK_TAG = "\n<script src=\"" + serverBase + "/js/jsHook.js\" type=\"text/javascript\"></script>"
	if len(tokenBuf) > 0 && tokenBuf[0] == '<' {
		token := string(tokenBuf)
		tagName := getTagName(token)
		if attrrs := tagAttrToProxy[tagName]; len(attrrs) != 0 {
			for _, attr := range attrrs {
				attrPattern := regexp.MustCompile("(" + attr + "=\\\")([^\\\"]+)(\\\")")
				if match := attrPattern.FindStringSubmatch(token); len(match) == 4 {
					rawUrl := match[2]
					proxiedUrl := GetProxiedUrl(rawUrl, url)
					token = attrPattern.ReplaceAllString(token, "${1}"+proxiedUrl+"${3}")
				}
			}
			tokenBuf = []byte(token)
		} else if tagName == "head" {
			tokenBuf = []byte(token + JS_HOOK_TAG)
		}

	}
	*htmlBuf = append(*htmlBuf, tokenBuf...)
}

func getTagName(token string) string {
	tagNameRegex := regexp.MustCompile("^<\\s*([a-zA-Z]+).*>$")
	submatch := tagNameRegex.FindStringSubmatch(token)
	if submatch != nil {
		return submatch[1]
	} else {
		return ""
	}
}

var tagAttrToProxy map[string][]string

func init() {
	tagAttrToProxy = map[string][]string{
		"a":      []string{"href"},
		"script": []string{"src"},
		"link":   []string{"href"},
		"base":   []string{"href"},
		"img":    []string{"src"},
		"meta":   []string{"content"},
		"form":   []string{"action"},
	}
}
