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
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/keys"
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
	b2, _ := ioutil.ReadFile("relay.history")
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
	b2, _ := ioutil.ReadFile("producer.history")
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
	} else if command == "keys" {
		digitalocean.ListKeys()
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
	} else if command == "relay" {
		dest := argMap["relay"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeRelay(argMap["producer"])
		ScpFile("setup.sh", dest)
		ScpFile("relay.sh", dest)
	} else if command == "fresh2" {
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
