package main

import (
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type LsData struct {
	Tip          Tip      `json:"tip"`
	Date         string   `json:"date"`
	SpecialFiles []string `json:"special_files"`
}

type Tip struct {
	Epoch int64  `json:"tip1"`
	Hash  string `json:"tip2"`
	Slot  int64  `json:"tip3"`
	Block int64  `json:"tip4"`
	Era   string `json:"tip5"`
}

func StartEyes() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		o, _ := exec.Command("ls", "-l", "/root/cardano-my-node/").Output()
		ls := LsData{}
		ls.Tip = Tip{}
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
