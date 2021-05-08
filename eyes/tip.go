package main

import (
	"os/exec"
)

func RunQueryTip() string {
	o, e := exec.Command("cardano-cli", "query", "tip", "--mainnet").CombinedOutput()
	if e != nil {
		return e.Error()
	}
	return string(o)
}
