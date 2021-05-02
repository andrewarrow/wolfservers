package runner

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/andrewarrow/wolfservers/files"
)

var PrivMap, PubMap map[string]string

func WriteOutJit(name string) {
	data := []byte(PrivMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit", data, 0600)
	data = []byte(PubMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit.pub", data, 0644)
}

func HotExec(name, ip, command string) string {
	WriteOutJit(name)
	env := "CARDANO_NODE_SOCKET_PATH=/root/cardano-my-node/db/socket"
	fullCommand := fmt.Sprintf("%s sudo -E %s", env, command)
	fmt.Println(fullCommand)
	o, e := exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, fullCommand).CombinedOutput()
	fmt.Println(e)
	return strings.TrimSpace(string(o))
}
