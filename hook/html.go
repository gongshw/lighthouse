package hook

import (
	"regexp"
)

const (
	JS_HOOK_TAG = "\n<script src=\"/js/jsHook.js\" type=\"text/javascript\"></script>"
)

func ParseHtml(html string, url string) string {
	base := GetBaseFromUrl(url)
	i := 0
	var htmlBuf []byte
	var tokenBuf []byte
	htmlLength := len(html)
	for {
		if i >= htmlLength {
			// end of html, flush tokenBuf
			if len(tokenBuf) > 0 {
				flushToken(&htmlBuf, tokenBuf, base)
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
			flushToken(&htmlBuf, tokenBuf, base)
			tokenBuf = tokenBuf[:0]
		} else {
			tokenBuf = append(tokenBuf, b)
			i++
			if b == '>' {
				flushToken(&htmlBuf, tokenBuf, base)
				tokenBuf = tokenBuf[:0]
			}
		}
	}
	return string(htmlBuf)
}

func flushToken(htmlBuf *[]byte, tokenBuf []byte, base string) {
	if len(tokenBuf) > 0 && tokenBuf[0] == '<' {
		if token := string(tokenBuf); needProxy(token) != "" {
			fullUrlRegex := regexp.MustCompile("(https?:\\/\\/([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?)")
			nonSchemaUrlRegex := regexp.MustCompile("(\\/\\/([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?)")
			absoluteUrlRegex := regexp.MustCompile("(\"\\s*)(\\/([\\/\\w \\.-]*)*\\/?)")
			if fullUrlRegex.MatchString(token) {
				token = fullUrlRegex.ReplaceAllString(token, "/proxy?$1")
			} else if nonSchemaUrlRegex.MatchString(token) {
				token = fullUrlRegex.ReplaceAllString(token, "/proxy?http:$1")
			} else if absoluteUrlRegex.MatchString(token) {
				token = absoluteUrlRegex.ReplaceAllString(token, "$1/proxy?"+base+"$2")
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
