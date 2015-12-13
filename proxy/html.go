package proxy

import (
	"github.com/gongshw/lighthouse/conf"
	"golang.org/x/net/html"
	"io"
	"strings"
)

var tagAttrToProxy map[string](map[string]bool)

var JS_HOOK_TAG string

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
		case html.StartTagToken:
			lastTag = flushTagToken(&newHtml, z, url)
		case html.SelfClosingTagToken:
			flushTagToken(&newHtml, z, url)
		}
		if tt != html.StartTagToken {
			lastTag = ""
		}
	}
}

func flushTagToken(htmlBuf *[]byte, tz *html.Tokenizer, url string) string {
	*htmlBuf = append(*htmlBuf, '<')
	tagName, hasAttr := tz.TagName()
	*htmlBuf = append(*htmlBuf, tagName...)
	if hasAttr {
		for {
			attrKey, attrValue, hasMore := tz.TagAttr()
			*htmlBuf = append(*htmlBuf, ' ')
			*htmlBuf = append(*htmlBuf, attrKey...)
			*htmlBuf = append(*htmlBuf, '=', '"')
			if tagAttrToProxy[string(tagName)][string(attrKey)] {
				urlInAttr := string(attrValue)
				*htmlBuf = append(*htmlBuf, []byte(GetProxiedUrl(urlInAttr, url))...)
			} else {
				*htmlBuf = append(*htmlBuf, attrValue...)
			}
			*htmlBuf = append(*htmlBuf, '"')
			if !hasMore {
				break
			}
		}
	}
	*htmlBuf = append(*htmlBuf, '>')
	if string(tagName) == "head" {
		*htmlBuf = append(*htmlBuf, []byte(getJsHookTag())...)
	}
	return string(tagName)
}

func getJsHookTag() string {
	if JS_HOOK_TAG == "" {
		serverBase := conf.CONFIG.ServerBaseUrl
		JS_HOOK_TAG = "\n<script src=\"" + serverBase + "/js/jsHook.js\" type=\"text/javascript\"></script>"
	}
	return JS_HOOK_TAG
}

func init() {
	tagAttrToProxyTemp := map[string][]string{
		"a":      []string{"href"},
		"script": []string{"src"},
		"link":   []string{"href"},
		"base":   []string{"href"},
		"img":    []string{"src"},
		"meta":   []string{"content"},
		"form":   []string{"action"},
	}
	tagAttrToProxy = make(map[string](map[string]bool))
	for tag, attrs := range tagAttrToProxyTemp {
		tagAttrToProxy[tag] = make(map[string]bool)
		for _, attr := range attrs {
			tagAttrToProxy[tag][attr] = true
		}
	}
}

