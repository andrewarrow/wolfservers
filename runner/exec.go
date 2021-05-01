package runner

import (
	"io/ioutil"

	"github.com/andrewarrow/wolfservers/files"
)

var PrivMap, PubMap map[string]string

func WriteOutJit(name string) {
	data := []byte(PrivMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit", data, 0600)
	data = []byte(PubMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit.pub", data, 0644)
}

func HotExec(command string) {
}
