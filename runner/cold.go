package runner

import (
	"fmt"
	"os/exec"

	"github.com/andrewarrow/wolfservers/files"
)

func ScpFileToCold(name, filename, dest string) {
	WriteOutJit(name)
	out, err := exec.Command("ssh", "-i", files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+dest,
		fmt.Sprintf("sudo cp /root/cardano-my-node/%s /home/aa/; sudo chown aa:aa /home/aa/%s", filename, filename)).CombinedOutput()
	fmt.Println(string(out))
	out, err = exec.Command("scp", "-i", files.UserHomeDir()+"/.ssh/wolf-jit", "aa@"+dest+":"+filename, ".").CombinedOutput()
	fmt.Println(string(out), err)
}
