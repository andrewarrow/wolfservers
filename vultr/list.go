package vultr

import (
	"context"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

func ListServers(pat string, ip2wolf map[string]string) {

	client, ctx := VultrClient(pat)
	listOptions := &govultr.ListOptions{PerPage: 100}
	i, _, _ := client.Instance.List(ctx, listOptions)
	for _, v := range i {
		wolfName := ip2wolf[v.MainIP]
		display.DisplayServer(wolfName, v.ID, "VULTR", v.Label, v.MainIP)
		//fmt.Println(v.ID, v.MainIP)
	}
}
func ListProducerIps(pat string) []string {

	list := []string{}
	client, ctx := VultrClient(pat)
	listOptions := &govultr.ListOptions{PerPage: 100}
	i, _, _ := client.Instance.List(ctx, listOptions)
	for _, v := range i {
		if v.Label != "producer" {
			continue
		}
		list = append(list, v.MainIP)
	}
	return list
}

func VultrClient(pat string) (*govultr.Client, context.Context) {
	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: pat})
	client := govultr.NewClient(oauth2.NewClient(ctx, ts))
	return client, ctx
}
