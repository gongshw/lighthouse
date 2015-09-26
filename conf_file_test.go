package main

import (
	"github.com/gongshw/lighthouse/conf"
	"testing"
)

func TestLighthouseJsonFile(t *testing.T) {
	err := conf.LoadConfig("./lighthouse.json.sample")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("lighthouse.json pass!")
	}
}
