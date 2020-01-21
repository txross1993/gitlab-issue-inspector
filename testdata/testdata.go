package testdata

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func HelperReadTestData(t *testing.T, filname string) []byte {
	_, dir, _, _ := runtime.Caller(1)
	path := filepath.Join(path.Dir(dir), "../testdata/", filname)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
