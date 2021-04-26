package keys

import (
	"crypto/rand"
	"fmt"
	"os/exec"

	"github.com/andrewarrow/wolfservers/files"
)

func WolfName() string {
	b := make([]byte, 16)
	rand.Read(b)
	name := fmt.Sprintf("wolf-%X", b[4:6])
	return name
}

func MakeEd() {
	out, err := exec.Command("ssh-keygen", "-o", "-a", "100", "-t", "ed25519",
		"-f", files.UserHomeDir()+"/.ssh/"+WolfName(), "-C", "wolfservers").Output()
	fmt.Println(string(out), err)
}
