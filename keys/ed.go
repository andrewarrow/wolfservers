package keys

import (
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
		if !sqlite.NameExists(name) {
			break
		}
		time.Sleep(time.Second)
		name = WolfName("wolf")
	}
	exec.Command("ssh-keygen", "-o", "-a", "100", "-t", "ed25519",
		"-f", files.UserHomeDir()+"/.ssh/"+name, "-C", "wolfservers").Output()

	file1 := files.UserHomeDir() + "/.ssh/" + name
	file2 := files.UserHomeDir() + "/.ssh/" + name + ".pub"
	b, _ := ioutil.ReadFile(file1)
	privKey := string(b)
	b, _ = ioutil.ReadFile(file2)
	pubKey := string(b)
	sqlite.InsertRow(name, provider, privKey, pubKey)
	os.Remove(file1)
	os.Remove(file2)
	return name, pubKey
}

func UpdateRowForEds(name string) {
	b, _ := ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name)
	privKey := string(b)
	b, _ = ioutil.ReadFile(files.UserHomeDir() + "/.ssh/" + name + ".pub")
	pubKey := string(b)
	sqlite.UpdateRow(name, privKey, pubKey)
}
