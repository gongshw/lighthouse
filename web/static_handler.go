package web

import (
	"fmt"
	"github.com/gongshw/lighthouse/bindata"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	WELCOME_PAGE = "index.html"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	var filePath string
	if indexOfQuestionMark := strings.Index(uri, "?"); indexOfQuestionMark > -1 {
		filePath = uri[1:indexOfQuestionMark]
	} else {
		filePath = uri[1:]
	}
	if filePath == "" {
		filePath = WELCOME_PAGE
	}
	if data, err := bindata.Asset(filePath); err == nil {
		var ctype string
		ctype = mime.TypeByExtension(filepath.Ext(filePath))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("404 page not found: %s", uri)))
	}
}
