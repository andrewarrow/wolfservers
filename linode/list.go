package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func ListServers(ip2wolf map[string]string) {

	client, ctx := LinodeClient()
	options := linodego.ListOptions{}

	i, _ := client.ListInstances(ctx, &options)
	for _, v := range i {
		wolfName := ip2wolf[fmt.Sprintf("%v", v.IPv4[0])]
		fmt.Printf("%s %s [LINODE] %s %s\n",
			wolfName,
			display.LeftAligned(v.ID, 10),
			display.LeftAligned(v.Label, 20),
			display.LeftAligned(fmt.Sprintf("ssh aa@%v", v.IPv4[0]), 40))
	}
}

func LinodeClient() (*linodego.Client, context.Context) {
	apiKey := os.Getenv("LINODE_PAT")

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

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
