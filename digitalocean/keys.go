package digitalocean

import (
	"encoding/json"
	"fmt"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/andrewarrow/wolfservers/network"
)

type Key struct {
	ID          int
	Fingerprint string
	Name        string
}
type Keys struct {
	Keys []Key `json:"ssh_keys"`
}

func ListKeys() {
	// ssh-keygen -l -E md5 -f  ~/.ssh/id_rsa
	jsonString := network.DoGet("v2/account/keys?per_page=100")
	//fmt.Println(jsonString)
	var keys Keys
	json.Unmarshal([]byte(jsonString), &keys)
	for _, s := range keys.Keys {
		fmt.Println(display.LeftAligned(s.ID, 5),
			display.LeftAligned(s.Fingerprint, 50),
			display.LeftAligned(s.Name, 10))
	}
}
