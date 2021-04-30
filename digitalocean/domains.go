package digitalocean

import (
	"encoding/json"
	"fmt"

	"github.com/andrewarrow/wolfservers/network"
)

// /v2/domains/$DOMAIN_NAME/records

func ListDomainRecords(domain string) []string {
	list := []string{}
	jsonString := network.DoGet(fmt.Sprintf("v2/domains/%s/records?per_page=100", domain))
	fmt.Println(jsonString)
	return list
}

type RecordThing struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
	Ttl  int    `json:"ttl"`
}

func AddRecord(domain, ip, name string) {
	ckt := RecordThing{}
	ckt.Name = name
	ckt.Type = "A"
	ckt.Data = ip
	ckt.Ttl = 86400
	asBytes, _ := json.Marshal(ckt)
	network.DoPost(fmt.Sprintf("v2/domains/%s/records", domain), asBytes)
}
