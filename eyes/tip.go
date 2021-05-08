package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

func RunQueryTip() string {
	o, e := exec.Command("cardano-cli", "query", "tip", "--mainnet").CombinedOutput()
	if e != nil {
		return e.Error()
	}
	return string(o)
}

func RunPaymentAmount() int64 {
	b, _ := ioutil.ReadFile("/root/cardano-my-node/payment.addr")
	o, _ := exec.Command("cardano-cli", "query", "utxo", "--address",
		string(b), "--mainnet").Output()
	for _, line := range strings.Split(string(o), "\n") {
		tokens := strings.Split(line, " ")
		fmt.Println("tokens", strings.Join(tokens, "|"))
		if len(tokens) == 15 {
			s := tokens[len(tokens)-2]
			i, _ := strconv.ParseInt(s, 10, 64)
			return i
		}

	}

	return 0
}
