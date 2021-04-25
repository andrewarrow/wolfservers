package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

func DigitalOcean() {
	fmt.Println("")
	fmt.Printf("%s\n", "List running droplets")
	fmt.Println("")
	ListDroplets()
}

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

func ListDroplets() {

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	fmt.Println("auth ok")

	ctx := context.TODO()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, _ := client.Droplets.List(ctx, opt)

	for _, droplet := range droplets {
		fmt.Printf("%v\n", droplet.Networks.V4[0].IPAddress)
	}
}
