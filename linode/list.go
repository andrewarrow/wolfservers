package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func ListServers() {

	client, ctx := LinodeClient()
	options := linodego.ListOptions{}

	res, err := client.ListInstances(ctx, &options)
	fmt.Println(res, err)
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
