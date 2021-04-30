package keys

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/sqlite"
)

func MakeEd(provider string) (string, string) {
	name := WolfName("wolf")
	for {
		filename := fmt.Sprintf("%s/.ssh/%s", files.UserHomeDir(), name)
		_, err := os.Stat(filename)
		if err != nil {
			break
		}
		time.Sleep(time.Second)
		name = WolfName("wolf")
	}
	exec.Command("ssh-keygen", "-o", "-a", "100", "-t", "ed25519",
		"-f", files.UserHomeDir()+"/.ssh/"+name, "-C", "wolfservers").Output()

	b, _ := ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name)
	privKey := string(b)
	b, _ = ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name + ".pub")
	pubKey := string(b)
	sqlite.InsertRow(name, provider, privKey, pubKey)
	return name, pubKey
}

func UpdateRowForEds(name string) {
	b, _ := ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name)
	privKey := string(b)
	b, _ = ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name + ".pub")
	pubKey := string(b)
	sqlite.UpdateRow(name, privKey, pubKey)
}
