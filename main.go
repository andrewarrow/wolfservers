package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
	"github.com/andrewarrow/wolfservers/keys"
	"github.com/andrewarrow/wolfservers/linode"
	"github.com/andrewarrow/wolfservers/sqlite"
	"github.com/andrewarrow/wolfservers/vultr"
	"github.com/justincampbell/timeago"
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

		fmt.Println("")
		if argMap["keys"] == "true" {
			vips := vultr.ListProducerIps()
			lips := linode.ListProducerIps()
			dips := digitalocean.ListProducerIps()

			ips := append(vips, lips...)
			ips = append(ips, dips...)
			for _, ip := range ips {
				out := SshAsUserRunOneThing("aa", ip2name[ip], ip)
				tokens := strings.Split(out, " ")
				month := tokens[5]
				day := tokens[6]
				hoursMins := tokens[7]
				ts, _ := time.Parse("Jan 2, 2006 15:04",
					fmt.Sprintf("%s %s, 2021 %s", month, day, hoursMins))
				fmt.Println(ip2name[ip], timeago.FromDuration(time.Since(ts)))

			}
		}
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
	} else if command == "producer" {
		// https://www.coincashew.com/coins/overview-ada/guide-how-to-build-a-haskell-stakepool-node
		dest := argMap["producer"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("scripts/node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeProducer(argMap["relay"])
		ScpFile(ip2name[dest], "setup.sh", dest)
		ScpFile(ip2name[dest], "producer.sh", dest)
		fmt.Println("1. ssh as root")
		fmt.Println("2. run setup.sh")
		fmt.Println("3. . .bashrc")
		fmt.Println("4. run producer.sh")
	} else if command == "relay" {
		dest := argMap["relay"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("scripts/node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeRelay(argMap["producer"])
		ScpFile(ip2name[dest], "setup.sh", dest)
		ScpFile(ip2name[dest], "relay.sh", dest)
		fmt.Println("1. ssh as root")
		fmt.Println("2. run setup.sh")
		fmt.Println("3. . .bashrc")
		fmt.Println("4. run relay.sh")
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
		// size := "s-4vcpu-8gb"
		size := "s-1vcpu-2gb"
		keys := digitalocean.ListKeyFingerprints()
		key := keys[0]
		digitalocean.CreateDroplet("producer", size, key)
		digitalocean.CreateDroplet("relay", size, key)

	} else if command == "node-keys" {
		//keys.MakeNode("wolf-C0B5")
	} else if command == "domains-do" {
		digitalocean.ListDomainRecords("wolfschedule.com")
	} else if command == "add-a-record" {
		ip := argMap["ip"]
		name := argMap["name"]
		digitalocean.AddRecord("wolfschedule.com", ip, name)
	} else if command == "danger-do" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(id)
	} else if command == "danger-linode" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		linode.RemoveServer(id)
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
		if argMap["provider"] == "" {
			return
		}
		provider := argMap["provider"]

		if provider == "linode" {
			ids := linode.ListKeys()
			for _, id := range ids {
				linode.DeleteSshKey(id)
			}
			name, pubKey := keys.MakeEd("LINODE")
			linode.CreateSshKey(name, strings.TrimSpace(pubKey))
		} else if provider == "do" {
			ids := digitalocean.ListKeys()
			for _, id := range ids {
				digitalocean.DeleteKey(id)
			}
			name, pubKey := keys.MakeEd("DO")
			digitalocean.CreateKey(name, pubKey)
		}
	} else if command == "help" {
		PrintHelp()
	}
}
