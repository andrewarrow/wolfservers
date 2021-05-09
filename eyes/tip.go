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
	fmt.Println("1111111111")
	b, _ := ioutil.ReadFile("/root/cardano-my-node/payment.addr")
	o, e := exec.Command("cardano-cli", "query", "utxo", "--address",
		strings.TrimSpace(string(b)), "--mainnet").CombinedOutput()
	fmt.Println(e, string(o))
	sum := int64(0)
	for _, line := range strings.Split(string(o), "\n") {
		tokens := strings.Split(line, " ")
		fmt.Println("tokens", strings.Join(tokens, "|"))
		if len(tokens) == 15 {
			s := tokens[len(tokens)-2]
			i, _ := strconv.ParseInt(s, 10, 64)
			sum += i
		}

	}

	return sum
}
