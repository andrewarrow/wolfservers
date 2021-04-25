package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
)

func PrintHelp() {
	fmt.Println("")
	fmt.Println("  wolfservers help         # this menu")
	fmt.Println("  wolfservers ls           # list servers")
	fmt.Println("  wolfservers make         # make new one --size=slug")
	fmt.Println("  wolfservers danger       # --ID=id")
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
	} else if command == "make" {
		if argMap["size"] == "" {
			digitalocean.ListSizes()
			return
		}
		digitalocean.CreateDroplet(argMap["size"])
	} else if command == "danger" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(id)
	} else if command == "help" {
		PrintHelp()
	}
}
