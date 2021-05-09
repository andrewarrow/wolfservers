package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andrewarrow/wolfservers/algorand"
	"github.com/andrewarrow/wolfservers/args"
	"github.com/andrewarrow/wolfservers/digitalocean"
	"github.com/andrewarrow/wolfservers/display"
	"github.com/andrewarrow/wolfservers/files"
	"github.com/andrewarrow/wolfservers/keys"
	"github.com/andrewarrow/wolfservers/linode"
	"github.com/andrewarrow/wolfservers/network"
	"github.com/andrewarrow/wolfservers/runner"
	"github.com/andrewarrow/wolfservers/sqlite"
	"github.com/andrewarrow/wolfservers/vultr"
	"github.com/justincampbell/timeago"
	touchid "github.com/lox/go-touchid"
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

func ProducerIps(pats map[string]string) []string {
	vips := vultr.ListProducerIps(pats["vultr"])
	lips := linode.ListProducerIps(pats["linode"])
	dips := digitalocean.ListProducerIps(pats["do"])

	ips := append(vips, lips...)
	ips = append(ips, dips...)
	return ips
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
	pats := sqlite.LoadPats()
	//ip2id := sqlite.MakeIpToId(db)
	runner.PrivMap, runner.PubMap = sqlite.SshKeysAsMap(db)

	tokens := strings.Split(command, ".")
	if len(tokens) == 4 {
		name := ip2name[command]
		SshAsUser("aa", name, command)
		return
	}

	if command == "ls" {
		digitalocean.ListDroplets(pats["do"], ip2name)
		vultr.ListServers(pats["vultr"], ip2name)
		linode.ListServers(pats["linode"], ip2name)

		fmt.Println("")
		if argMap["keys"] == "true" {
			sum := int64(0)
			for _, ip := range ProducerIps(pats) {
				jsonString := network.DoIpGet(ip)
				var ls LsDataHolder
				json.Unmarshal([]byte(jsonString), &ls)
				tokens := strings.Split(ls.M.Date, " ")
				month := tokens[0]
				day := tokens[1]
				hoursMins := tokens[2]
				if len(tokens) == 5 {
					month = tokens[0]
					day = tokens[2]
					hoursMins = tokens[3]
				}

				ts, _ := time.Parse("Jan 2, 2006 15:04",
					fmt.Sprintf("%s %s, 2021 %s", month, day, hoursMins))
				ago := timeago.FromDuration(time.Since(ts))
				fmt.Printf("%s %s %d %0.2f %d %d\n", ip2name[ip],
					display.LeftAligned(ago, 20),
					ls.M.Tip.Epoch, float64(ls.M.Amount)/1000000.0, ls.M.Tip.Slot,
					ls.M.Tip.Block)
				sum += ls.M.Amount

				/*
					// Apr 27 19:16 kes.vkey

					fmt.Printf("%s (%d,%d,%d) Key Evolving Signature Age(%s) \n",
						ip2name[ip], rsd.Tip.Epoch, rsd.Tip.Block, rsd.Tip.Slot, timeago.FromDuration(time.Since(ts)))
					fmt.Printf("    |------> %s\n", rsd.SpecialFiles)
				*/
			}
			fmt.Printf("Total %0.2f\n", float64(sum)/1000000.0)
		}
	} else if command == "phases" {
		fmt.Println("")
		fmt.Println(" 1. Seed     : new VM from provider, ip address known")
		fmt.Println(" 2. Growing  : cardano-node installed, sync started")
		fmt.Println(" 3. Growing  : sync complete, now generate block-producer keys")
		fmt.Println(" 4. Growing  : also setup payment, stake keys, send real ADA")
		fmt.Println(" 5. Growing  : build 1st tx, copy tx.raw to cold, sign, run hot")
		fmt.Println(" 6. Growing  : register your stake pool with DNS A record")
		fmt.Println(" 7. Growing  : pool.cert,deleg.cert to hot")
		fmt.Println(" 8. Growing  : build 2nd tx, copy tx.raw to cold, sign, run hot")
		fmt.Println(" 9. Exciting : find your data on block explorers pooltool.io")
		fmt.Println("10. Next     : configure your topology files")
		fmt.Println("11. Next     : wait four hours!")
		fmt.Println("12. Then     : fetch your relay node buddies")
		fmt.Println("13. MakeSure : must see the Processed TX number increasing")
		fmt.Println("14. AndThen  : wait an epoch or two and then...")
		fmt.Println("15. OMG      : checking stake pool rewards!")
		fmt.Println("")
	} else if command == "next" {
		ip := argMap["ip"]
		name := ip2name[ip]
		nodeKeySize := sqlite.NodeKeysQuery(name)
		paymentAddr := sqlite.PaymentKeysQuery(name)
		fmt.Println(nodeKeySize, paymentAddr)
		ioutil.WriteFile("payment.addr", []byte(paymentAddr), 0755)
		ScpFileToHot(name, "payment.addr", ip)
		os.Remove("payment.addr")
		//if paymentAddr == "" {
		//	keys.MakePayment(name)
		//}
	} else if command == "tx-delegate" {
		ip := argMap["ip"]
		name := ip2name[ip]
		ps, ss := sqlite.PaymentAndStakeSigning(name)
		ioutil.WriteFile("payment.skey", []byte(ps), 0755)
		ioutil.WriteFile("stake.skey", []byte(ss), 0755)
		sqlite.CreateNodeKeysOnDisk(name)
		keys.SignTxDelegate()
		ScpFileToHot(name, "tx.signed", ip)
		os.Remove("tx.raw")
		os.Remove("tx.signed")
		os.Remove("payment.skey")
		os.Remove("stake.skey")
		os.Remove("node.counter")
		os.Remove("node.skey")
		os.Remove("node.vkey")
		result := runner.HotExec(name, ip, "cardano-cli transaction submit --tx-file /root/cardano-my-node/tx.signed --mainnet")
		fmt.Println(result)
	} else if command == "tx" {
		ip := argMap["ip"]
		name := ip2name[ip]
		ps, ss := sqlite.PaymentAndStakeSigning(name)
		ioutil.WriteFile("payment.skey", []byte(ps), 0755)
		ioutil.WriteFile("stake.skey", []byte(ss), 0755)
		keys.SignTx()
		ScpFileToHot(name, "tx.signed", ip)
		os.Remove("tx.raw")
		os.Remove("tx.signed")
		os.Remove("payment.skey")
		os.Remove("stake.skey")
		result := runner.HotExec(name, ip, "cardano-cli transaction submit --tx-file /root/cardano-my-node/tx.signed --mainnet")
		fmt.Println(result)
	} else if command == "deleg.cert" {
		ip := argMap["ip"]
		name := ip2name[ip]
		sv := sqlite.PaymentStakeV(name)
		ioutil.WriteFile("stake.vkey", []byte(sv), 0755)
		sqlite.CreateNodeKeysOnDisk(name)
		keys.Delegation()
		ScpFileToHot(name, "deleg.cert", ip)
		os.Remove("stake.vkey")
		os.Remove("node.counter")
		os.Remove("node.skey")
		os.Remove("node.vkey")
		os.Remove("deleg.cert")
	} else if command == "pool.cert" {
		ip := argMap["ip"]
		name := ip2name[ip]
		hash := "5318dc07f229acbace49e666124b528c99a36763857f67a2f379d428166577fa"
		code := "3A81"
		sv := sqlite.PaymentStakeV(name)
		ioutil.WriteFile("stake.vkey", []byte(sv), 0755)
		runner.ScpFileToCold(name, "vrf.vkey", ip)
		sqlite.CreateNodeKeysOnDisk(name)
		// 9997821299
		// 9997.82
		keys.StakePoolRegCert(9000, 2, hash, code)
		ScpFileToHot(name, "pool.cert", ip)
		os.Remove("stake.vkey")
		os.Remove("node.counter")
		os.Remove("node.skey")
		os.Remove("node.vkey")
		os.Remove("pool.cert")
	} else if command == "stake.cert" {
		ip := argMap["ip"]
		name := ip2name[ip]
		sv := sqlite.PaymentStakeV(name)
		ioutil.WriteFile("stake.vkey", []byte(sv), 0755)
		keys.MakeStakeCert()
		ScpFileToHot(name, "stake.cert", ip)
		os.Remove("stake.vkey")
		os.Remove("stake.cert")
	} else if command == "keys" {
		digitalocean.ListKeys()
	} else if command == "cold" {
		ip := argMap["ip"]
		name := ip2name[ip]
		runner.ScpFileToCold(name, "tx.raw", ip)
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
		ScpFileToHot(name, "node.cert", ip)
		// 6. delete local kes.vkey, node.skey
		os.Remove("kes.vkey")
		os.Remove("node.cert")
		os.Remove("node.counter")
		os.Remove("node.skey")
		os.Remove("node.vkey")
	} else if command == "algorand" {
		algorand.Query()
	} else if command == "add-pat" {
		provider := argMap["provider"]
		pat := argMap["pat"]
		sqlite.InsertPat(provider, pat)
	} else if command == "show-oath" {
		// oathtool --totp -b ''
		sqlite.ShowOaths()
	} else if command == "add-oath" {
		name := argMap["name"]
		seed := argMap["seed"]
		username := argMap["username"]
		password := argMap["password"]
		sqlite.InsertOath(name, seed, username, password)
	} else if command == "comments" {
		in := argMap["in"]
		files.RemoveComments(in)
	} else if command == "touch" {

		ok, err := touchid.Authenticate("access llamas")
		if err != nil {
			log.Fatal(err)
		}

		if ok {
			log.Printf("Authenticated")
		} else {
			log.Fatal("Failed to authenticate")
		}
	} else if command == "deploy" {
		for _, ip := range ProducerIps(pats) {
			name := ip2name[ip]
			ScpFileToAa(name, "eyes.next", ip)
		}
	} else if command == "temp" {
		ip := argMap["ip"]
		name := ip2name[ip]
		ScpFileToHot(name, "scripts/stake.register", ip)
		ScpFileToHot(name, "scripts/delegate.pool", ip)
	} else if command == "poolMetaData" {
		code := argMap["code"]
		GenPoolMetaData(code)
	} else if command == "ready-params" {
		ip := argMap["ip"]
		name := ip2name[ip]
		o := runner.HotExec(name, ip, "cardano-cli query protocol-parameters --mainnet --out-file /root/cardano-my-node/params.json")
		fmt.Println(o)
	} else if command == "sqlite" {
		sqlite.List()
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
	} else if command == "producer" {
		// https://www.coincashew.com/coins/overview-ada/guide-how-to-build-a-haskell-stakepool-node
		dest := argMap["producer"]
		PrepDest(dest)
		MakeProducer(argMap["relay"])
		ScpFile(ip2name[dest], "scripts/cardano.setup", dest)
		ScpFile(ip2name[dest], "scripts/stake.register", dest)
		ScpFile(ip2name[dest], "scripts/delegate.pool", dest)
		ScpFile(ip2name[dest], "producer.sh", dest)
		fmt.Println("1. ssh as root")
		fmt.Println("2. run cardano.setup")
		fmt.Println("3. . .bashrc")
		fmt.Println("4. run producer.sh")
	} else if command == "relay" {
		dest := argMap["relay"]
		PrepDest(dest)
		MakeRelay(argMap["producer"])
		ScpFile(ip2name[dest], "scripts/cardano.setup", dest)
		ScpFile(ip2name[dest], "scripts/stake.register", dest)
		ScpFile(ip2name[dest], "scripts/delegate.pool", dest)
		ScpFile(ip2name[dest], "relay.sh", dest)
		fmt.Println("1. ssh as root")
		fmt.Println("2. run cardano.setup")
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
		linode.CreateServer(pats["linode"], "producer")
		linode.CreateServer(pats["linode"], "relay")
	} else if command == "fresh2vultr" {
		if argMap["sure"] == "" {
			return
		}
		vultr.CreateServer(pats["vultr"], "producer")
		vultr.CreateServer(pats["vultr"], "relay")
	} else if command == "fresh2do" {
		if argMap["sure"] == "" {
			return
		}
		// make 2 droplets, name one producer one relay, wait for their ips
		// size := "s-4vcpu-8gb"
		size := "s-1vcpu-2gb"
		keys := digitalocean.ListKeyFingerprints()
		key := keys[0]
		digitalocean.CreateDroplet(pats["do"], "producer", size, key)
		digitalocean.CreateDroplet(pats["do"], "relay", size, key)

	} else if command == "node-keys" {
		ip := argMap["ip"]
		name := ip2name[ip]
		keys.MakeNode(name)
		keys.MakePayment(name)
		ScpFileToHot(name, "payment.addr", ip)
	} else if command == "domains-do" {
		digitalocean.ListDomainRecords(pats["do"], "wolfschedule.com")
	} else if command == "add-a-record" {
		ip := argMap["ip"]
		name := argMap["name"]
		digitalocean.AddRecord(pats["do"], "wolfschedule.com", ip, name)
	} else if command == "danger-do" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		digitalocean.RemoveDroplet(pats["do"], id)
	} else if command == "danger-vultr" {
		if argMap["ID"] == "" {
			return
		}
		vultr.RemoveServer(pats["vultr"], argMap["ID"])
	} else if command == "danger-linode" {
		if argMap["ID"] == "" {
			return
		}
		id, _ := strconv.Atoi(argMap["ID"])
		linode.RemoveServer(pats["linode"], id)
	} else if command == "ed255" {
		if argMap["provider"] == "" {
			return
		}
		provider := argMap["provider"]

		if provider == "linode" {
			ids := linode.ListKeys(pats["linode"])
			for _, id := range ids {
				linode.DeleteSshKey(pats["linode"], id)
			}
			name, pubKey := keys.MakeEd("LINODE")
			linode.CreateSshKey(pats["linode"], name, strings.TrimSpace(pubKey))
		} else if provider == "vultr" {
			ids := vultr.ListKeys(pats["vultr"])
			for _, id := range ids {
				vultr.DeleteKey(pats["vultr"], id)
			}
			name, pubKey := keys.MakeEd("VULTR")
			vultr.CreateKey(pats["vultr"], name, strings.TrimSpace(pubKey))
		} else if provider == "do" {
			ids := digitalocean.ListKeys()
			for _, id := range ids {
				digitalocean.DeleteKey(id)
			}
			name, pubKey := keys.MakeEd("DO")
			digitalocean.CreateKey(name, strings.TrimSpace(pubKey))
		}
	} else if command == "help" {
		PrintHelp()
	}
}
