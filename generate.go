package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/andrewarrow/wolfservers/files"
)

func unescape(s string) template.HTML {
	return template.HTML(s)
}

type Replacer struct {
	RelayIP        string
	LtLt           string
	ProducerIP     string
	StartKesPeriod string
}

func MakeRelay(ip string) {
	b2, _ := ioutil.ReadFile("scripts/relay.history")
	blob := string(b2)
	t := template.Must(template.New("relay").
		Funcs(template.FuncMap{"unescape": unescape}).
		Parse(blob))
	var buff bytes.Buffer
	r := Replacer{}
	r.ProducerIP = ip
	r.LtLt = "<<"
	t.Execute(&buff, r)
	ioutil.WriteFile("relay.sh", buff.Bytes(), 0755)
}
func MakeProducer(ip string) {
	b2, _ := ioutil.ReadFile("scripts/producer.history")
	blob := string(b2)
	t := template.Must(template.New("producer").
		Funcs(template.FuncMap{"unescape": unescape}).
		Parse(blob))
	var buff bytes.Buffer
	r := Replacer{}
	r.RelayIP = ip
	r.LtLt = "<<"
	t.Execute(&buff, r)
	ioutil.WriteFile("producer.sh", buff.Bytes(), 0755)
}
func MakeAirGap(StartKesPeriod string) {
	b2, _ := ioutil.ReadFile("scripts/airgapped.keys")
	blob := string(b2)
	t := template.Must(template.New("thing").
		Parse(blob))
	var buff bytes.Buffer
	r := Replacer{}
	r.StartKesPeriod = StartKesPeriod
	t.Execute(&buff, r)
	ioutil.WriteFile("airgapped/airgapped.sh", buff.Bytes(), 0755)
}

func PrepDest(dest string) {
	out, _ := exec.Command("ssh-keyscan", "-H", dest).Output()
	f, _ := os.OpenFile(files.UserHomeDir()+"/.ssh/known_hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(out))
}

var privMap, pubMap map[string]string

func WriteOutJit(name string) {
	data := []byte(privMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit", data, 0600)
	data = []byte(pubMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit.pub", data, 0644)
}
func SshAsUser(user, name, ip string) {
	WriteOutJit(name)
	/*
		  rwx | 7 | Read, write and execute  |
		| rw- | 6 | Read, write              |
		| r-x | 5 | Read, and execute        |
		| r-- | 4 | Read,                    |
		| -wx | 3 | Write and execute        |
		| -w- | 2 | Write                    |
		| --x | 1 | Execute                  |
		| --- | 0 | no permissions           |
		| rwx------  | 0700 | User  |
		| ---rwx---  | 0070 | Group |

			-rw-------  1 andrewarrow  staff   399 Apr 26 19:28 wolf-91F4
		  -rw-r--r--  1 andrewarrow  staff    81 Apr 26 19:28 wolf-91F4.pub
	*/

	fmt.Println("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", user+"@"+ip)

	cmd := exec.Command("ssh", "-i", files.UserHomeDir()+"/.ssh/wolf-jit", user+"@"+ip)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")
	FancyPtyAndTerm(cmd)
}

func ScpFile(name, file, dest string) {
	WriteOutJit(name)
	out, err := exec.Command("scp", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", file, "root@"+dest+":").Output()
	fmt.Println(string(out), err)
}
func ScpFileFromRemote(orig string) {
	out, err := exec.Command("mkdir", "airgapped").Output()
	out, err = exec.Command("scp", "aa@"+orig+":kes*", "airgapped/").Output()
	fmt.Println(string(out), err)
}
func ScpFileToHot(file, dest string) {
	out, err := exec.Command("scp", file, "aa@"+dest+":").Output()
	fmt.Println(string(out), err)
	tokens := strings.Split(file, "/")
	filename := tokens[1]
	out, _ = exec.Command("ssh", "aa@"+dest,
		fmt.Sprintf("sudo cp %s /root/cardano-my-node; rm %s;", filename, filename)).CombinedOutput()
	fmt.Println(string(out))
}
func ScpFileToNodeHome(file, dest string) string {
	out, err := exec.Command("scp", file, "aa@"+dest+":").Output()
	fmt.Println(string(out), err)
	//tokens := strings.Split(file, "/")
	//out, err = exec.Command("ssh", "aa@"+dest, fmt.Sprintf("'sudo mv %s /root/cardano-my-node/'", tokens[1])).CombinedOutput()
	out, _ = exec.Command("ssh", "aa@"+dest, "sudo cp producer.keys /root/cardano-my-node; rm producer.keys; sudo chmod +x /root/cardano-my-node/producer.keys; sudo /root/cardano-my-node/producer.keys; sudo cp /root/cardano-my-node/kes.* /home/aa").CombinedOutput()
	tokens := strings.Split(string(out), "\n")
	for _, line := range tokens {
		if strings.HasPrefix(line, "startKesPeriod") {
			tokens = strings.Split(line, ":")
			return tokens[1]
		}
	}
	return ""
}