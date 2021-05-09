package digitalocean

import (
	"encoding/json"
	"fmt"

	"github.com/andrewarrow/wolfservers/network"
)

// /v2/domains/$DOMAIN_NAME/records
type DomainRecord struct {
	Name string
	Data string
}
type DomainRecordsThing struct {
	DomainRecords []DomainRecord `json:"domain_records"`
}

func ListDomainRecords(pat, domain string) []string {
	list := []string{}
	jsonString := network.DoGet(pat, fmt.Sprintf("v2/domains/%s/records?per_page=100", domain))
	var drt DomainRecordsThing
	json.Unmarshal([]byte(jsonString), &drt)
	for _, thing := range drt.DomainRecords {
		fmt.Printf("%30s.wolfschedule.com %30s\n", thing.Name, thing.Data)
	}
	return list
}

type RecordThing struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
	Ttl  int    `json:"ttl"`
}

func AddRecord(pat, domain, ip, name string) {
	ckt := RecordThing{}
	ckt.Name = name
	ckt.Type = "A"
	ckt.Data = ip
	ckt.Ttl = 86400
	asBytes, _ := json.Marshal(ckt)
	network.DoPost(pat, fmt.Sprintf("v2/domains/%s/records", domain), asBytes)
}
