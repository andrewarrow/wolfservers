package digitalocean

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/andrewarrow/wolfservers/network"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var (
	pat = os.Getenv("DO_PAT")
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

/*
func CreateDroplet() {
	dropletName := "super-cool-droplet"

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "SFO3",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-14-04-x64",
		},
	}

	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("err: %s\n\n", err)

	}
	fmt.Printf(newDroplet.PublicIPv4())
}*/

func ListSizes() {
	jsonString := network.DoGet("v2/sizes")
	var sizes DropletSizes
	json.Unmarshal([]byte(jsonString), &sizes)
	for _, s := range sizes.Sizes {
		fmt.Println(s)
	}
}
func ListDroplets() {

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, _ := client.Droplets.List(ctx, opt)

	for _, droplet := range droplets {
		fmt.Printf("%s %v\n", droplet.Name, droplet.Networks.V4[1].IPAddress)
	}
}
