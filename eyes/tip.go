package main

import (
	"os/exec"
)

func RunQueryTip() string {
	o, _ := exec.Command("cardano-cli query tip --mainnet").Output()
	return string(o)
}
