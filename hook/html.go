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
	fullUrlRegex      = regexp.MustCompile("(https?):\\/\\/(([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?)")
	nonSchemaUrlRegex = regexp.MustCompile("\\/\\/([\\da-z\\.-]+\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?)")
	absoluteUrlRegex  = regexp.MustCompile("(\"\\s*)(\\/([\\/\\w \\.-]*)*\\/?)")
)

func flushToken(htmlBuf *[]byte, tokenBuf []byte, url string) {
	var serverBase string = conf.CONFIG.ServerBaseUrl
	var JS_HOOK_TAG = "\n<script src=\"" + serverBase + "/js/jsHook.js\" type=\"text/javascript\"></script>"
	if len(tokenBuf) > 0 && tokenBuf[0] == '<' {
		if token := string(tokenBuf); needProxy(token) != "" {
			if fullUrlRegex.MatchString(token) {
				token = fullUrlRegex.ReplaceAllString(token, serverBase+"/proxy/${1}/${2}")
			} else if nonSchemaUrlRegex.MatchString(token) {
				token = nonSchemaUrlRegex.ReplaceAllString(token, serverBase+"/proxy/http/$1")
			} else if absoluteUrlRegex.MatchString(token) {
				protocal, host := ParseBaseUrl(url)
				rplcUrl := "${1}" + serverBase + "/proxy/" + protocal + "/" + host + "/" + "$2"
				token = absoluteUrlRegex.ReplaceAllString(token, rplcUrl)
			}
			tokenBuf = []byte(token)
		} else if getTagName(token) == "head" {
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

func needProxy(token string) string {
	tagToProxy := map[string]string{
		"a":      "href",
		"script": "src",
		"link":   "href",
		"base":   "href",
		"img":    "src",
		"meta":   "content",
		"form":   "action",
	}
	return tagToProxy[getTagName(token)]
}
