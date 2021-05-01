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
	"github.com/andrewarrow/wolfservers/runner"
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
	fmt.Println("  wolfservers phases       # how to do what when")
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
	runner.PrivMap, runner.PubMap = sqlite.SshKeysAsMap(db)

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
				rsd := SshAsUserRunOneThing(ip2name[ip], ip)
				// Apr 27 19:16 kes.vkey
				tokens := strings.Split(rsd.Date, " ")
				month := tokens[0]
				day := tokens[1]
				hoursMins := tokens[2]
				ts, _ := time.Parse("Jan 2, 2006 15:04",
					fmt.Sprintf("%s %s, 2021 %s", month, day, hoursMins))

				fmt.Printf("%s EPOCH(%d) Key Evolving Signature Age(%s) \n",
					ip2name[ip], rsd.Tip.Epoch, timeago.FromDuration(time.Since(ts)))
				fmt.Printf("    |------> %s\n", rsd.SpecialFiles)

			}
		}
	} else if command == "phases" {
		fmt.Println("")
		fmt.Println(" 1. Seed     : new VM from provider, ip address known")
		fmt.Println(" 2. Growing  : cardano-node installed, sync started")
		fmt.Println(" 3. Growing  : sync complete, now generate block-producer keys")
		fmt.Println(" 4. Growing  : also setup payment and stake keys")
		fmt.Println(" 5. Growing  : then register your stake address")
		fmt.Println(" 6. Growing  : register your stake pool with DNS A record")
		fmt.Println(" 7. Growing  : pool.cert,deleg.cert to hot")
		fmt.Println(" 8. Growing  : build 1st tx, copy tx.raw to cold")
		fmt.Println(" 9. Growing  : make and copy tx.signed to hot. Execute")
		fmt.Println("10. Exciting : find your data on block explorers pooltool.io")
		fmt.Println("11. Next     : configure your topology files")
		fmt.Println("12. Next     : wait four hours!")
		fmt.Println("13. Then     : fetch your relay node buddies")
		fmt.Println("14. MakeSure : must see the Processed TX number increasing")
		fmt.Println("15. AndThen  : wait an epoch or two and then...")
		fmt.Println("16. OMG      : checking stake pool rewards!")
		fmt.Println("")
	} else if command == "keys" {
		digitalocean.ListKeys()
	} else if command == "hot" {
		ip := argMap["ip"]
		RunHot(ip2name[ip], ip)
	} else if command == "issue-op-cert" {
		ip := argMap["ip"]
		name := ip2name[ip]
		// 1. download kes.vkey
		CatKesV(name, ip)
		// 2. use node.skey from sqlite
		sqlite.CreateNodeKeysOnDisk(name)
		// 3. get startKesPeriod from hot
		startKesPeriod := RunHot(name, ip)
		// 4. keys.IssueOpCert(startKesPeriod)
		keys.IssueOpCert(startKesPeriod)
		// 5. upload node.cert to hot
		ScpFileToHot("node.cert", ip)
		// 6. delete local kes.vkey, node.skey
		os.Remove("kes.vkey")
		os.Remove("node.cert")
		os.Remove("node.counter")
		os.Remove("node.skey")
		os.Remove("node.vkey")
	} else if command == "ready-params" {
		ip := argMap["ip"]
		name := ip2name[ip]
		runner.HotExec(name, ip, "cardano-cli query protocol-parameters --mainnet --out-file params.json")
	} else if command == "sqlite" {
		sqlite.List()
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
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
		//keys.MakePayment("wolf-C0B5")
		//ScpFileToHot("payment.addr", ip)
		//keys.MakeStakeCert(name)
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
