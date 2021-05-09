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

func ListKeys() []int {
	list := []int{}
	jsonString := network.DoGet("", "v2/account/keys?per_page=100")
	//fmt.Println(jsonString)
	var keys Keys
	json.Unmarshal([]byte(jsonString), &keys)
	for _, s := range keys.Keys {
		fmt.Println(display.LeftAligned(s.ID, 5),
			display.LeftAligned(s.Fingerprint, 50),
			display.LeftAligned(s.Name, 10))
		list = append(list, s.ID)
	}
	return list
}

type CreateKeyThing struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func CreateKey(name, pubKey string) {
	ckt := CreateKeyThing{}
	ckt.Name = name
	ckt.PublicKey = pubKey
	asBytes, _ := json.Marshal(ckt)
	network.DoPost("", "v2/account/keys", asBytes)
}

func DeleteKey(id int) {
	network.DoDelete(fmt.Sprintf("v2/account/keys/%d", id))
}
func ListKeyFingerprints() []string {
	list := []string{}
	jsonString := network.DoGet("", "v2/account/keys?per_page=100")
	//fmt.Println(jsonString)
	var keys Keys
	json.Unmarshal([]byte(jsonString), &keys)
	for _, s := range keys.Keys {
		fmt.Println(display.LeftAligned(s.ID, 5),
			display.LeftAligned(s.Fingerprint, 50),
			display.LeftAligned(s.Name, 10))
		list = append(list, s.Fingerprint)
	}
	return list
}
