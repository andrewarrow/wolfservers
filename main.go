package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
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
	fmt.Println("")
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
	} else if command == "wolf" {
		fmt.Println("groupadd ssh-users")
		fmt.Println("useradd -c 'get in sync' -m -d /home/wolf -s /bin/bash -G sudo,ssh-users wolf")
		fmt.Println("rsync --archive --chown=wolf:wolf ~/.ssh /home/wolf")

		text := `apt update
apt upgrade
apt install -y build-essential libssl-dev pkg-config nload jq python3-pip tcptraceroute chrony
cd /usr/bin; wget http://www.vdberg.org/~richard/tcpping; chmod 755 tcpping; cd
curl -LO https://github.com/BurntSushi/ripgrep/releases/download/11.0.2/ripgrep_11.0.2_amd64.deb
dpkg -i ripgrep_11.0.2_amd64.deb; rm ripgrep_11.0.2_amd64.deb
`

		fmt.Println(text)
	} else if command == "make" {
		if argMap["size"] == "" {
			digitalocean.ListSizes()
			return
		}
		key := os.Getenv("DO_PRINT")
		digitalocean.CreateDroplet(argMap["size"], key)
	} else if command == "danger" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(id)
	} else if command == "ed255" {
		out, err := exec.Command("ssh-keygen", "-o", "-a", "100", "-t", "ed25519",
			"-f", "/Users/andrewarrow/.ssh/id_ed25519", "-C", "wolfservers").Output()
		fmt.Println(string(out), err)
	} else if command == "help" {
		PrintHelp()
	}
}
