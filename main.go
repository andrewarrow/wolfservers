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
	fmt.Println("  wolfservers images       # list images")
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
	} else if command == "images" {
		digitalocean.ListImages(1)
		digitalocean.ListImages(2)
	} else if command == "wolf2" {
		// https://www.coincashew.com/coins/overview-ada/guide-how-to-build-a-haskell-stakepool-node
		fmt.Println("sudo apt-get update -y")
		fmt.Println("sudo apt-get upgrade -y")
		fmt.Println("sudo apt-get install git jq bc make automake rsync htop curl build-essential pkg-config libffi-dev libgmp-dev libssl-dev libtinfo-dev libsystemd-dev zlib1g-dev make g++ wget libncursesw5 libtool autoconf -y")
		fmt.Println("")
		l := `mkdir $HOME/git
cd $HOME/git
git clone https://github.com/input-output-hk/libsodium
cd libsodium
git checkout 66f017f1
./autogen.sh
./configure
make
sudo make install`
		fmt.Println(l)
		fmt.Println("sudo apt-get -y install pkg-config libgmp-dev libssl-dev libtinfo-dev libsystemd-dev zlib1g-dev build-essential curl libgmp-dev libffi-dev libncurses-dev libtinfo5")
		fmt.Println("curl --proto '=https' --tlsv1.2 -sSf https://get-ghcup.haskell.org | sh")
		l = `cd $HOME
source .bashrc
ghcup upgrade
ghcup install cabal 3.4.0.0
ghcup set cabal 3.4.0.0`
		fmt.Println(l)
		l = `ghcup install ghc 8.10.4
ghcup set ghc 8.10.4`
		fmt.Println(l)
		fmt.Println("")
		fmt.Println("")

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
