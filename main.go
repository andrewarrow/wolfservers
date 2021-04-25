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
	fmt.Println("  wolfservers make         # make new one --size=slug --key=key")
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
		fmt.Println("useradd -c “get in sync” -m -d /home/wolf -s /bin/bash -G sudo,ssh-users wolf")
		fmt.Println("rsync --archive --chown=wolf:wolf ~/.ssh /home/wolf")

		text := `apt update
apt upgrade
apt install -y build-essential libssl-dev
apt install pkg-config
apt install nload
apt install jq
apt install python3-pip
apt-get install tcptraceroute
apt-get install chrony`

		fmt.Println(text)
	} else if command == "make" {
		if argMap["size"] == "" || argMap["key"] == "" {
			digitalocean.ListSizes()
			return
		}
		digitalocean.CreateDroplet(argMap["size"], argMap["key"])
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
