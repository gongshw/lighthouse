package ssl

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateTempCrtAndKey(t *testing.T) {
	crtPath, keyPath, err := CreateTempCrtAndKey("localhost")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("crtPath: %s;keyPath: %s", crtPath, keyPath)
		crtData, _ := ioutil.ReadFile(crtPath)
		keyData, _ := ioutil.ReadFile(keyPath)
		fmt.Println(string(crtData))
		fmt.Println(string(keyData))
		os.Remove(crtPath)
		os.Remove(keyPath)
	}
}
