package main

import (
	"encoding/json"
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

func QueryStakeAddress() int64 {
	b, _ := ioutil.ReadFile("/root/cardano-my-node/stake.addr")
	o, e := exec.Command("cardano-cli", "query", "stake-address-info",
		"--address", strings.TrimSpace(string(b)), "--mainnet").CombinedOutput()
	/*
		[
		    {
		        "address": "stake1uxju8kl78g40zm5fnffwz8mgrl22yp80c4dvzdtpw9m8rpcavekux",
		        "rewardAccountBalance": 0,
		        "delegation": "pool1fqnluegrkns2jj49vgvn8vvn5u8y4m2m3xptx4wngvf4cnwa2fk"
		    }
		]

	*/
	var list []interface{}
	json.Unmarshal(o, &list)
	if len(list) == 0 {
		fmt.Println(string(o), e)
		return 0
	}
	m := list[0].(map[string]interface{})
	balance := m["rewardAccountBalance"].(float64)
	return int64(balance)
}
