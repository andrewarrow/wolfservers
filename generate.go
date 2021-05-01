package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/runner"
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

func PrepDest(dest string) {
	out, _ := exec.Command("ssh-keyscan", "-H", dest).Output()
	f, _ := os.OpenFile(files.UserHomeDir()+"/.ssh/known_hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(out))
}

func RunHot(name, ip string) int {
	nodeHome := "/root/cardano-my-node"
	command := fmt.Sprintf("sudo cat %s/mainnet-shelley-genesis.json | jq -r '.slotsPerKESPeriod'", nodeHome)
	runner.WriteOutJit(name)
	o, _ := exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, command).Output()
	slotsPerKESPeriod := strings.TrimSpace(string(o))

	env := "CARDANO_NODE_SOCKET_PATH=/root/cardano-my-node/db/socket"
	command = fmt.Sprintf("%s sudo -E cardano-cli query tip --mainnet | jq -r '.slot'", env)
	o, _ = exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, command).Output()
	slot := strings.TrimSpace(string(o))

	slotsPerKESPeriodInt, _ := strconv.Atoi(slotsPerKESPeriod)
	slotInt, _ := strconv.Atoi(slot)
	startKesPeriod := slotInt / slotsPerKESPeriodInt
	return startKesPeriod
}
func SshAsUserRunOneThing(name, ip string) string {
	runner.WriteOutJit(name)
	// cardano-cli query tip --mainnet
	o, _ := exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, "CARDANO_NODE_SOCKET_PATH=/root/cardano-my-node/db/socket sudo -E cardano-cli query tip --mainnet").Output()
	fmt.Println(string(o))
	o, _ = exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, "sudo ls -l /root/cardano-my-node/kes.vkey").Output()
	return string(o)
}
func CatKesV(name, ip string) {
	runner.WriteOutJit(name)
	out, _ := exec.Command("ssh", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+ip, "sudo cat /root/cardano-my-node/kes.vkey").Output()
	data := strings.TrimSpace(string(out))
	ioutil.WriteFile("kes.vkey", []byte(data), 0755)
}

func SshAsUser(user, name, ip string) {
	runner.WriteOutJit(name)
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
	runner.WriteOutJit(name)
	out, err := exec.Command("scp", "-i",
		files.UserHomeDir()+"/.ssh/wolf-jit", file, "root@"+dest+":").Output()
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
