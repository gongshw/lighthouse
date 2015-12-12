package proxy

import (
	"testing"
)

func TestParseBaseUrl(t *testing.T) {
	var url, protocal, host string

	url = "http://gongsw.com"
	protocal, host = ParseBaseUrl(url)
	if protocal != "http" || host != "gongsw.com" {
		t.Errorf("fail to run ParseBaseUrl on url %s", url)
	}

	url = "http://gongsw.com/"
	protocal, host = ParseBaseUrl(url)
	if protocal != "http" || host != "gongsw.com" {
		t.Errorf("fail to run ParseBaseUrl on url %s", url)
	}

	url = "http://gongsw.com/test/path/"
	protocal, host = ParseBaseUrl(url)
	if protocal != "http" || host != "gongsw.com" {
		t.Errorf("fail to run ParseBaseUrl on url %s", url)
	}

	url = "http//gongsw.com/test/file.txt"
	protocal, host = ParseBaseUrl(url)
	if protocal != "" || host != "" {
		t.Errorf("fail to run ParseBaseUrl on a wrong url %s", url)
	}
}

func TestParseUrl(t *testing.T) {
	var url, protocal, uri string
	url = "http://gongsw.com"
	protocal, uri = ParseUrl(url)
	if protocal != "http" || uri != "gongsw.com/" {
		t.Errorf("fail to run ParseUrl on url %s", url)
	}

	url = "http://gongsw.com/"
	protocal, uri = ParseUrl(url)
	if protocal != "http" || uri != "gongsw.com/" {
		t.Errorf("fail to run ParseUrl on url %s", url)
	}

	url = "http://gongsw.com/test/path/"
	protocal, uri = ParseUrl(url)
	if protocal != "http" || uri != "gongsw.com/test/path/" {
		t.Errorf("fail to run ParseUrl on url %s", url)
	}

	url = "http//gongsw.com/test/file.txt"
	protocal, uri = ParseUrl(url)
	if protocal != "" || uri != "" {
		t.Errorf("fail to run ParseUrl on a wrong url %s", url)
	}
}

func TestGetProxiedUrl(t *testing.T) {
	var base, url string
	base = "http://gongshw.com/path/to/file"

	url = "/"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "/path"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/path" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "//gongshw.com"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "//gongshw.com/"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "//gongshw.com/path/"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/path/" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "//gongshw.com/path/to/file"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/path/to/file" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "http://gongshw.com/path/to/file"
	if GetProxiedUrl(url, base) != "/proxy/http/gongshw.com/path/to/file" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "https://gongshw.com/path/to/file"
	if GetProxiedUrl(url, base) != "/proxy/https/gongshw.com/path/to/file" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}

	url = "path/to/file"
	if GetProxiedUrl(url, base) != "path/to/file" {
		t.Errorf("fail to run GetProxiedUrl on (%s,%s)", url, base)
	}
}
