package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/keys"
	"github.com/andrewarrow/wolfservers/linode"
	"github.com/andrewarrow/wolfservers/sqlite"
	"github.com/andrewarrow/wolfservers/vultr"
	"github.com/creack/pty"
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
	fmt.Println("  wolfservers sqlite       # list data in sqlite db")
	fmt.Println("  wolfservers ssh          # --ip=")
	fmt.Println("")
}

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

func SshAsUser(user, name, ip string) {
	data := []byte(privMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit", data, 0600)
	data = []byte(pubMap[name])
	ioutil.WriteFile(files.UserHomeDir()+"/.ssh/wolf-jit.pub", data, 0644)
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
	test(cmd)
}
func test(c *exec.Cmd) error {

	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	defer func() { _ = ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	defer func() { signal.Stop(ch); close(ch) }()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func ScpFile(name, file, dest string) {
	out, err := exec.Command("scp", "-i",
		files.UserHomeDir()+"/.ssh/"+name, file, "root@"+dest+":").Output()
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
func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 1 {
		PrintHelp()
		return
	}
	command := os.Args[1]
	argMap := args.ToMap()

	db := sqlite.OpenTheDB()
	defer db.Close()
	ip2name := sqlite.MakeIpMap(db)
	ip2id := sqlite.MakeIpToId(db)
	privMap, pubMap = sqlite.SshKeysAsMap(db)

	if command == "ls" {
		digitalocean.ListDroplets(ip2name)
		vultr.ListServers(ip2name)
		linode.ListServers(ip2name)
	} else if command == "keys" {
		digitalocean.ListKeys()
	} else if command == "sqlite" {
		sqlite.List()
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
	} else if command == "wolfit" {
		ip := argMap["producer"]
		startKesPeriod := ScpFileToNodeHome("scripts/producer.keys", ip)
		ScpFileFromRemote(ip)
		MakeAirGap(startKesPeriod)
	} else if command == "wolfit2" {
		ip := argMap["producer"]
		ScpFileToHot("airgapped/node.cert", ip)
	} else if command == "relay" {
		dest := argMap["relay"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("scripts/node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeRelay(argMap["producer"])
		ScpFile(ip2name[dest], "setup.sh", dest)
		ScpFile(ip2name[dest], "relay.sh", dest)
	} else if command == "update-ips" {
		producer := argMap["producer"]
		relay := argMap["relay"]
		name := argMap["name"]
		sqlite.UpdateIps(name, producer, relay)
	} else if command == "ssh" {
		ip := argMap["ip"]
		name := ip2name[ip]
		user := "aa"
		if argMap["root"] == "true" {
			user = "root"
		}
		SshAsUser(user, name, ip)

	} else if command == "update-ed" {
		name := argMap["name"]
		keys.UpdateRowForEds(name)
	} else if command == "update-ids" {
		producer := argMap["producer"]
		relay := argMap["relay"]
		name := argMap["name"]
		sqlite.UpdateIds(name, producer, relay)
	} else if command == "fresh2linode" {
		if argMap["sure"] == "" {
			return
		}
		linode.CreateServer("producer")
		linode.CreateServer("relay")
	} else if command == "fresh2vultr" {
		if argMap["sure"] == "" {
			return
		}
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
		ScpFile(ip2name[dest], "setup.sh", dest)
		ScpFile(ip2name[dest], "producer.sh", dest)
	} else if command == "danger-do" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(id)
	} else if command == "danger" {
		if argMap["name"] == "" {
			return
		}
		for k, v := range ip2name {
			if v == argMap["name"] {
				sid := ip2id[k]
				id, _ := strconv.Atoi(sid)
				linode.RemoveServer(id)
			}
		}
	} else if command == "ed255" {
		name, pubKey := keys.MakeEd("LINODE")
		linode.CreateSshKey(name, strings.TrimSpace(pubKey))
	} else if command == "help" {
		PrintHelp()
	}
}
