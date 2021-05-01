package keys

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/andrewarrow/wolfservers/sqlite"
)

func IssueOpCert(startKesPeriod int) {
	fmt.Println("IssueOpCert", startKesPeriod)
	o, err := exec.Command("cardano-cli", "node", "issue-op-cert",
		"--kes-verification-key-file", "kes.vkey",
		"--cold-signing-key-file", "node.skey",
		"--operational-certificate-issue-counter", "node.counter",
		"--kes-period", fmt.Sprintf("%d", startKesPeriod),
		"--out-file", "node.cert").CombinedOutput()
	fmt.Println(string(o), err)
}
func MakeNode(name string) {
	exec.Command("cardano-cli", "node", "key-gen", "--cold-verification-key-file",
		"node.vkey", "--cold-signing-key-file", "node.skey",
		"--operational-certificate-issue-counter", "node.counter").Output()

	b, _ := ioutil.ReadFile("node.vkey")
	v := string(b)
	b, _ = ioutil.ReadFile("node.skey")
	s := string(b)
	b, _ = ioutil.ReadFile("node.counter")
	c := string(b)
	sqlite.InsertNodeRow(name, v, s, c)
	os.Remove("node.vkey")
	os.Remove("node.skey")
	os.Remove("node.counter")
}
func ToTokens(s string) []string {
	return strings.Split(s, " ")
}
func MakePayment(name string) {
	cmd := "cardano-cli"
	tokens := ToTokens("address key-gen --verification-key-file payment.vkey --signing-key-file payment.skey")
	exec.Command(cmd, tokens...).Output()
	tokens = ToTokens("stake-address key-gen --verification-key-file stake.vkey --signing-key-file stake.skey")
	exec.Command(cmd, tokens...).Output()
	tokens = ToTokens("stake-address build --stake-verification-key-file stake.vkey --out-file stake.addr --mainnet")
	exec.Command(cmd, tokens...).Output()
	tokens = ToTokens("address build --payment-verification-key-file payment.vkey --stake-verification-key-file stake.vkey --out-file payment.addr --mainnet")
	exec.Command(cmd, tokens...).Output()

	b, _ := ioutil.ReadFile("payment.vkey")
	pv := string(b)
	b, _ = ioutil.ReadFile("payment.skey")
	ps := string(b)
	b, _ = ioutil.ReadFile("stake.vkey")
	sv := string(b)
	b, _ = ioutil.ReadFile("stake.skey")
	ss := string(b)
	b, _ = ioutil.ReadFile("stake.addr")
	sa := string(b)
	b, _ = ioutil.ReadFile("payment.addr")
	pa := string(b)
	sqlite.InsertPaymentRow(name, pv, ps, sv, ss, sa, pa)

	os.Remove("payment.vkey")
	os.Remove("payment.skey")
	os.Remove("stake.vkey")
	os.Remove("stake.skey")
	os.Remove("stake.addr")
	os.Remove("payment.addr")
}

func Step1() {
	/*
						cardano-cli node key-gen-KES \
						    --verification-key-file kes.vkey \
						    --signing-key-file kes.skey


		  cardano-cli node key-gen \
				    --cold-verification-key-file node.vkey \
				    --cold-signing-key-file node.skey \
				    --operational-certificate-issue-counter node.counter

			cardano-cli node issue-op-cert \
		    --kes-verification-key-file kes.vkey \
		    --cold-signing-key-file node.skey \
		    --operational-certificate-issue-counter node.counter \
		    --kes-period <startKesPeriod> \
		    --out-file node.cert
	*/
}

/*

S = sign    == private
V = verify  == public

stake pool cold key (node.cert)
stake pool hot key (kes.skey)        Key Evolving Signature
stake pool VRF key (vrf.skey)        Verifiable random function

you will need to regenerate the KES key every 90 days.

Determine the number of slots per KES period
kesPeriod by dividing the slot tip number by the slotsPerKESPeriod
operational certificate

1. Copy kes.vkey to your cold environment.
2. Copy node.cert to your hot environment.


Payment keys are used to send and receive payments and stake keys are used to manage stake delegations.

payment.skey & payment.vkey
stake.skey & stake.vkey

stake.addr

Build a payment address for the payment key payment.vkey which will delegate to the stake address, stake.vkey

payment.addr

4. Copy payment.addr to your hot environment.
5. Create a certificate, stake.cert, using the stake.vkey

stake.cert

6. Copy stake.cert to your hot environment.

7. Copy tx.raw to your cold environment.
8. Copy tx.signed to your hot environment.

pool.cert

Copy pool.cert to your hot environment.

deleg.cert

Copy deleg.cert to your hot environment.

Copy tx.raw to your cold environment.
Copy tx.signed to your hot environment.

Copy stakepoolid.txt to your hot environment.


./topologyUpdater.sh

Complete this section after four hours when your relay node IP is properly registered.

./relay-topology_pull.sh


*/
