package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/keys"
	"github.com/andrewarrow/wolfservers/linode"
	"github.com/andrewarrow/wolfservers/vultr"
)

func PrintHelp() {
	fmt.Println("")
	fmt.Println("  wolfservers help         # this menu")
	fmt.Println("  wolfservers ls           # list servers")
	fmt.Println("  wolfservers keys         # list ssh keys")
	fmt.Println("  wolfservers make         # make new one --size=slug")
	fmt.Println("  wolfservers danger       # --ID=id")
	fmt.Println("  wolfservers ed255        # new ed25519 key")
	fmt.Println("  wolfservers wolf         # user add for wolf user")
	fmt.Println("  wolfservers images       # list images")
	fmt.Println("  wolfservers tags         # list tags")
	fmt.Println("")
}

func unescape(s string) template.HTML {
	return template.HTML(s)
}

type Replacer struct {
	RelayIP    string
	LtLt       string
	ProducerIP string
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
func ScpFile(file, dest string) {
	out, err := exec.Command("scp", file, "root@"+dest+":").Output()
	fmt.Println(string(out), err)
}
func ScpFileFromRemote(file, orig string) {
	out, err := exec.Command("scp", "aa@"+orig+":kes*", ".").Output()
	fmt.Println(string(out), err)
}
func ScpFileToNodeHome(file, dest string) {
	out, err := exec.Command("scp", file, "aa@"+dest+":").Output()
	fmt.Println(string(out), err)
	//tokens := strings.Split(file, "/")
	//out, err = exec.Command("ssh", "aa@"+dest, fmt.Sprintf("'sudo mv %s /root/cardano-my-node/'", tokens[1])).CombinedOutput()
	out, _ = exec.Command("ssh", "aa@"+dest, "sudo cp producer.keys /root/cardano-my-node; rm producer.keys; sudo chmod +x /root/cardano-my-node/producer.keys; sudo /root/cardano-my-node/producer.keys; sudo cp /root/cardano-my-node/kes.* /home/aa").CombinedOutput()
	tokens := strings.Split(string(out), "\n")
	for _, line := range tokens {
		if strings.HasPrefix(line, "startKesPeriod") {
			tokens = strings.Split(line, ":")
			fmt.Println("startKesPeriod!", tokens[1])
		}
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 1 {
		PrintHelp()
		return
	}
	command := os.Args[1]
	argMap := args.ToMap()

	if command == "ls" {
		digitalocean.ListDroplets()
		vultr.ListServers()
		linode.ListServers()
	} else if command == "keys" {
		digitalocean.ListKeys()
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
	} else if command == "wolfit" {
		ip := argMap["producer"]
		ScpFileToNodeHome("scripts/producer.keys", ip)
		ScpFileFromRemote("", ip)
		// scripts/producer.keys
		//   kes.vkey
		//   kes.skey
		//   startKesPeriod: 215
	} else if command == "relay" {
		dest := argMap["relay"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("scripts/node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeRelay(argMap["producer"])
		ScpFile("setup.sh", dest)
		ScpFile("relay.sh", dest)
	} else if command == "fresh2linode" {
		linode.CreateServer("producer")
		linode.CreateServer("relay")
	} else if command == "fresh2vultr" {
		vultr.CreateServer("producer")
		vultr.CreateServer("relay")
	} else if command == "fresh2do" {
		if argMap["sure"] == "" {
			return
		}
		// make 2 droplets, name one producer one relay, wait for their ips
		size := "s-4vcpu-8gb"
		key := os.Getenv("DO_PRINT")
		digitalocean.CreateDroplet("producer", size, key)
		digitalocean.CreateDroplet("relay", size, key)

	} else if command == "producer" {
		// https://www.coincashew.com/coins/overview-ada/guide-how-to-build-a-haskell-stakepool-node
		dest := argMap["producer"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeProducer(argMap["relay"])
		ScpFile("setup.sh", dest)
		ScpFile("producer.sh", dest)
	} else if command == "danger" {
		if argMap["ID"] == "" {
			return
		}
		//if false { // TODO rethink how to prevent disaster
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(id)
		//}
	} else if command == "ed255" {
		keys.MakeEd()
	} else if command == "help" {
		PrintHelp()
	}
}
