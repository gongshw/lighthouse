package main

import (
	"github.com/gongshw/lighthouse/conf"
	"testing"
)

func TestLghthouseJsonFile(t *testing.T) {
	err := conf.LoadConfig("./lighthouse.json")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("lighthouse.json pass!")
	}
}
