package vultr

import (
	"context"
	"fmt"
	"os"

	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

func ListServers() {

	client, ctx := VultrClient()
	listOptions := &govultr.ListOptions{PerPage: 100}
	i, _, _ := client.Instance.List(ctx, listOptions)
	for _, v := range i {
		fmt.Println(v)
	}
}

func VultrClient() (*govultr.Client, context.Context) {
	apiKey := os.Getenv("VULTR_PAT")

	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: apiKey})
	client := govultr.NewClient(oauth2.NewClient(ctx, ts))
	return client, ctx
}
