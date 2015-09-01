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

var (
	styleTokenPattern = regexp.MustCompile("(style\\s*=\\s*\\\")([^\\\"]+)(\\\")")
)

func flushToken(htmlBuf *[]byte, tokenBuf []byte, url string) {
	var serverBase string = conf.CONFIG.ServerBaseUrl
	var JS_HOOK_TAG = "\n<script src=\"" + serverBase + "/js/jsHook.js\" type=\"text/javascript\"></script>"
	if len(tokenBuf) > 0 && tokenBuf[0] == '<' {
		token := string(tokenBuf)
		if tagName := getTagName(token); tagName != "" && tagName[0] != '/' {
			// a tag open

			if match := styleTokenPattern.FindStringSubmatch(token); len(match) == 4 {
				css := match[2]
				token = styleTokenPattern.ReplaceAllString(token, "${1}"+ParseCss(css, url)+"${3}")
			}

			if attrrs := tagAttrToProxy[tagName]; len(attrrs) != 0 {
				for _, attr := range attrrs {
					// replace all links in tag attr
					attrPattern := regexp.MustCompile("(" + attr + "=\\\")([^\\\"]+)(\\\")")
					if match := attrPattern.FindStringSubmatch(token); len(match) == 4 {
						rawUrl := match[2]
						proxiedUrl := GetProxiedUrl(rawUrl, url)
						token = attrPattern.ReplaceAllString(token, "${1}"+proxiedUrl+"${3}")
					}
				}
			}

			if tagName == "head" {
				// inject js
				token = token + JS_HOOK_TAG
			}

			tokenBuf = []byte(token)
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
