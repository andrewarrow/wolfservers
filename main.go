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

func MakeRelay() {
	b2, _ := ioutil.ReadFile("relay.history")
	blob := string(b2)
	t := template.Must(template.New("relay").
		Funcs(template.FuncMap{"unescape": unescape}).
		Parse(blob))
	var buff bytes.Buffer
	r := Replacer{}
	r.ProducerIP = "165.232.150.86"
	r.LtLt = "<<"
	t.Execute(&buff, r)
	ioutil.WriteFile("relay.sh", buff.Bytes(), 0755)
}
func MakeProducer() {
	b2, _ := ioutil.ReadFile("producer.history")
	blob := string(b2)
	t := template.Must(template.New("producer").
		Funcs(template.FuncMap{"unescape": unescape}).
		Parse(blob))
	var buff bytes.Buffer
	r := Replacer{}
	r.RelayIP = "144.126.222.70"
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
		if argMap["dest"] == "" {
			return
		}
		dest := argMap["dest"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeRelay()
		ScpFile("setup.sh", dest)
		ScpFile("relay.sh", dest)
	} else if command == "producer" {
		if argMap["dest"] == "" {
			return
		}
		// https://www.coincashew.com/coins/overview-ada/guide-how-to-build-a-haskell-stakepool-node
		dest := argMap["dest"]
		PrepDest(dest)
		b1, _ := ioutil.ReadFile("node.setup")
		ioutil.WriteFile("setup.sh", b1, 0755)
		MakeProducer()
		ScpFile("setup.sh", dest)
		ScpFile("producer.sh", dest)
		// apply tag
	} else if command == "wolf" {

		fmt.Println("useradd -c 'aa' -m -d /home/aa -s /bin/bash -G sudo aa")
		fmt.Println("rsync --archive --chown=aa:aa ~/.ssh /home/aa")
		fmt.Println("sudo bash -c 'echo \"aa ALL=(ALL) NOPASSWD:ALL\" >> /etc/sudoers.d/99_sudo_include_file'")

		text := `apt update
apt upgrade
apt install -y build-essential libssl-dev pkg-config nload jq python3-pip tcptraceroute chrony
cd /usr/bin; wget http://www.vdberg.org/~richard/tcpping; chmod 755 tcpping; cd
curl -LO https://github.com/BurntSushi/ripgrep/releases/download/11.0.2/ripgrep_11.0.2_amd64.deb
dpkg -i ripgrep_11.0.2_amd64.deb; rm ripgrep_11.0.2_amd64.deb
ufw disable
ufw default deny incoming; ufw default allow outgoing
ufw limit proto tcp from any to any port 22
ufw allow proto tcp from any to any port 3000
ufw enable
su - wolf
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
git clone https://github.com/input-output-hk/jormungandr
cd jormungandr/
git submodule update --init --recursive
cargo install --path jormungandr --force

git clone https://github.com/Chris-Graffagnino/Jormungandr-for-Newbs.git -b files-only --single-branch files

chmod +x ~/files/*.sh; chmod +x ~/files/env
cat ~/files/.bashrc > ~/.bashrc && cat ~/files/.bash_profile > ~/.bash_profile && cat ~/files/.bash_aliases > .bash_aliases

chmod 700 ~/.bashrc && chmod 700 ~/.bash_profile

source ~/.bash_profile

mkdir /home/wolf/storage

echo "export USERNAME='wolf'" >> ~/.bashrc
echo "export PUBLIC_IP_ADDR='161.35.232.233'" >> ~/.bashrc
echo "export REST_PORT='3001'" >> ~/.bashrc
echo "export REST_URL='http://127.0.0.1:3001/api'" >> ~/.bashrc
echo "export JORMUNGANDR_STORAGE_DIR='/home/wolf/storage'" >> ~/.bashrc
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
