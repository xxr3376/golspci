package lspci

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("testcase/output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	std, err := ioutil.ReadFile("testcase/output-std.json")
	if err != nil {
		t.Fatal(err)
	}
	stdData := make(map[string]map[string]string)
	err = json.Unmarshal(std, &stdData)
	if err != nil {
		t.Fatal(err)
	}

	data, err := parseLSPCI(f)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(data, stdData) {
		t.Fail()
	}
}
