package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type LsData struct {
	Tip          Tip      `json:"tip"`
	Date         string   `json:"date"`
	SpecialFiles []string `json:"special_files"`
	Amount       int64    `json:"amount"`
	Balance      int64    `json:"balance"`
}

type Tip struct {
	Epoch int64  `json:"epoch"`
	Hash  string `json:"hash"`
	Slot  int64  `json:"slot"`
	Block int64  `json:"block"`
	Era   string `json:"era"`
}

func StartEyes() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		jsonString := RunQueryTip()
		fmt.Println(jsonString)
		var tip Tip
		json.Unmarshal([]byte(jsonString), &tip)

		o, _ := exec.Command("ls", "-l", "/root/cardano-my-node/").Output()
		ls := LsData{}
		ls.Amount = RunPaymentAmount()
		ls.Balance = QueryStakeAddress()
		ls.Tip = tip
		ls.SpecialFiles = []string{}

		special := []string{"pool.cert", "params.json", "node.cert", "payment.addr", "stake.cert", "tx.raw", "tx.signed", "poolMetaData.json"}
		for _, line := range strings.Split(string(o), "\n") {
			if strings.Contains(line, "kes.vkey") {
				ls.Date = strings.TrimSpace(line[31:])
			}
			for _, s := range special {
				if strings.Contains(line, s) {
					ls.SpecialFiles = append(ls.SpecialFiles, s)
				}
			}
		}

		c.JSON(200, gin.H{"m": ls})
	})
	r.Run(":8081")
}

func main() {
	StartEyes()
}
