package hook

import (
	"github.com/gongshw/lighthouse/conf"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

func ParseHtml(r io.Reader, url string) ([]byte, error) {
	z := html.NewTokenizer(r)
	var newHtml []byte
	lastTag := ""
	for {
		tt := z.Next()
		rawHtmlBytes := z.Raw()
		switch tt {
		case html.ErrorToken:
			e := z.Err()
			if e.Error() == "EOF" {
				return newHtml, nil
			} else {
				return make([]byte, 0), z.Err()
			}
		case html.TextToken:
			rawHtml := strings.TrimSpace(string(rawHtmlBytes[:]))
			if len(rawHtml) > 0 && lastTag == "style" {
				newCss := ParseCss(rawHtml, url)
				newHtml = append(newHtml, []byte(newCss)...)
			} else {
				newHtml = append(newHtml, rawHtmlBytes...)
			}
		case html.DoctypeToken, html.CommentToken, html.EndTagToken:
			newHtml = append(newHtml, rawHtmlBytes...)
		case html.StartTagToken, html.SelfClosingTagToken:
			flushToken(&newHtml, z, url)
		}
		if tt == html.StartTagToken {
			lastTagByte, _ := z.TagName()
			lastTag = string(lastTagByte[:])
		} else {
			lastTag = ""
		}
	}
}

var (
	styleTokenPattern = regexp.MustCompile("(style\\s*=\\s*\\\")([^\\\"]+)(\\\")")
)

func flushToken(htmlBuf *[]byte, tz *html.Tokenizer, url string) {
	tokenRaw := tz.Raw()
	var serverBase string = conf.CONFIG.ServerBaseUrl
	var JS_HOOK_TAG = "\n<script src=\"" + serverBase + "/js/jsHook.js\" type=\"text/javascript\"></script>"
	if len(tokenRaw) > 0 && tokenRaw[0] == '<' {
		token := string(tokenRaw)
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

			tokenRaw = []byte(token)
		}
	}
	*htmlBuf = append(*htmlBuf, tokenRaw...)
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
