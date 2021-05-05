package linode

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func ListServers(pat string, ip2wolf map[string]string) {

	client, ctx := LinodeClient(pat)
	options := linodego.ListOptions{}

	i, _ := client.ListInstances(ctx, &options)
	for _, v := range i {
		wolfName := ip2wolf[fmt.Sprintf("%v", v.IPv4[0])]
		display.DisplayServer(wolfName, v.ID, "LINODE", v.Label, v.IPv4[0])
	}
}
func ListProducerIps(pat string) []string {

	list := []string{}
	client, ctx := LinodeClient(pat)
	options := linodego.ListOptions{}

	i, _ := client.ListInstances(ctx, &options)
	for _, v := range i {
		if v.Label != "producer" {
			continue
		}
		list = append(list, fmt.Sprintf("%v", v.IPv4[0]))
	}
	return list
}

func LinodeClient(pat string) (*linodego.Client, context.Context) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: pat})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)
	client.SetDebug(false)
	ctx := context.Background()
	return &client, ctx
}

func ListKeys(pat string) []int {

	list := []int{}
	client, ctx := LinodeClient(pat)
	options := linodego.ListOptions{}

	keys, _ := client.ListSSHKeys(ctx, &options)
	for _, v := range keys {
		list = append(list, v.ID)
	}
	return list
}
