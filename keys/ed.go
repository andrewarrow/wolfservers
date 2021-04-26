package keys

import (
	"fmt"
	"os/exec"

	"github.com/andrewarrow/wolfservers/files"
)

func MakeEd() {
	out, err := exec.Command("ssh-keygen", "-o", "-a", "100", "-t", "ed25519",
		"-f", files.UserHomeDir()+"/.ssh/"+WolfName("wolf"), "-C", "wolfservers").Output()
	fmt.Println(string(out), err)
}
