package web

import (
	"fmt"
	"github.com/gongshw/lighthouse/conf"
	"html/template"
	"log"
	"net/http"
)

func ShowError(w http.ResponseWriter, msg string, url string) {
	if t, err := template.ParseFiles(conf.CONFIG.StaicFileDir + "/error.html"); err == nil {
		t.Execute(w, map[string]string{"msg": msg, "url": url})
	} else {
		log.Println("can't find error.html in %s", conf.CONFIG.StaicFileDir)
		fmt.Fprintf(w, msg+": %s", url)
	}
}
