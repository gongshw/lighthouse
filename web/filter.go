package web

import (
	"bufio"
	"errors"
	"github.com/gongshw/lighthouse/conf"
	"github.com/gongshw/lighthouse/proxy"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	FILTER_MODE_WHITE_LIST = iota
	FILTER_MODE_BLACK_LIST = iota
)

const (
	BLOCK_ACTION_REDIRECT  = iota
	BLOCK_ACTION_SHOW_INFO = iota
)

var (
	filterInited = false

	patterns = []*regexp.Regexp{}

	filterMode = FILTER_MODE_BLACK_LIST
)

func UrlNeedProxy(url string) bool {
	if conf.CONFIG.FilterMode == "" {
		// no filter config
		return true
	}
	_, host := proxy.ParseBaseUrl(url)
	for _, p := range patterns {
		if p.MatchString(host) {
			return filterMode == FILTER_MODE_WHITE_LIST
		}
	}
	return filterMode == FILTER_MODE_BLACK_LIST
}

func InitFilter() error {
	if !filterInited {
		if conf.CONFIG.FilterMode == "" {
			// no filter config
			filterInited = true
			return nil
		}
		if strings.EqualFold(conf.CONFIG.FilterMode, "white") {
			filterMode = FILTER_MODE_WHITE_LIST
		} else if strings.EqualFold(conf.CONFIG.FilterMode, "black") {
			filterMode = FILTER_MODE_BLACK_LIST
		} else {
			return errors.New("can't resolve filter mode: " + conf.CONFIG.FilterMode)
		}
		filterListFile, _ := filepath.Abs(conf.CONFIG.FilterFile)
		if file, err := os.Open(filterListFile); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if re := strings.TrimSpace(scanner.Text()); len(re) > 0 {
					if pattern, err := regexp.Compile(re); err == nil {
						patterns = append(patterns, pattern)
					} else {
						return errors.New("wrong regexp: " + re)
					}
				} else {
					//ignore
					continue
				}

			}
			log.Println("load proxy filter file from " + filterListFile)
			filterInited = true
			return nil
		} else {
			return errors.New("read filter file fail! " + err.Error())
		}
	}
	return nil
}
