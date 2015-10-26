package config

import (
	"os"
	"testing"
)

var initCasesContext = `
process = --args=/bin/bash --cwd=/bin --terminal=true;--args=/bin/bash;--cwd=/bin
`

func TestConfig(t *testing.T) {
	f, err := os.Create("cases.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(initCasesContext)
	if err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove("cases.conf")

	if data := GetConfig("process", "cases.conf"); len(data) != 3 {
		t.Fatal("Get process err", data)
	} else if data[0] != "--args=/bin/bash --cwd=/bin --terminal=true" {
		t.Fatal("Get first params of process err")
	}
}
