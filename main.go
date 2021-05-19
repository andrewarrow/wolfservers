package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/runner"
	"github.com/andrewarrow/wolfservers/sqlite"
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
	fmt.Println("  wolfservers phases       # how to do what when")
	fmt.Println("")
}

var argMap map[string]string
var ip2name map[string]string
var pats map[string]string

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 1 {
		PrintHelp()
		return
	}
	command := os.Args[1]
	argMap = args.ToMap()

	db := sqlite.OpenTheDB()
	defer db.Close()
	ip2name = sqlite.MakeIpMap(db)
	pats = sqlite.LoadPats()
	runner.PrivMap, runner.PubMap = sqlite.SshKeysAsMap(db)

	tokens := strings.Split(command, ".")
	if len(tokens) == 4 {
		name := ip2name[command]
		SshAsUser("aa", name, command)
		return
	}

	if command == "ls" {
		MainList()
	} else if command == "add-a-record" {
		AddARecord()
	} else if command == "add-pat" {
	} else if command == "add-oath" {
		AddOath()
	} else if command == "address" {
		PaymentAddress()
	} else if command == "cold" {
	} else if command == "comments" {
		Comments()
	} else if command == "danger-do" {
		DangerDo()
	} else if command == "danger-vultr" {
		DangerVultr()
	} else if command == "danger-linode" {
		DangerLinode()
	} else if command == "deleg.cert" {
	} else if command == "deploy" {
		DeployEyes()
	} else if command == "domains-do" {
	} else if command == "ed255" {
	} else if command == "fresh2linode" {
	} else if command == "fresh2vultr" {
	} else if command == "fresh2do" {
	} else if command == "hot" {
	} else if command == "issue-op-cert" {
	} else if command == "images" {
	} else if command == "keys" {
	} else if command == "node-keys" {
	} else if command == "poolMetaData" {
	} else if command == "producer" {
	} else if command == "pool.cert" {
	} else if command == "phases" {
		Phases()
	} else if command == "relay" {
	} else if command == "ready-params" {
	} else if command == "ssh" {
		MainSsh()
	} else if command == "show-oath" {
		sqlite.ShowOaths()
	} else if command == "stake.cert" {
	} else if command == "sqlite" {
	} else if command == "touch" {
	} else if command == "tx-delegate" {
	} else if command == "tx" {
		RunTx()
	} else if command == "update-ed" {
	} else if command == "update-ids" {
	} else if command == "update-ips" {
		UpdateIps()
	} else if command == "help" {
		PrintHelp()
	}
}
