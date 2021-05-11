package digitalocean

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/andrewarrow/wolfservers/network"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
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

func ListSizes() {
	jsonString := network.DoGet("", "v2/sizes?per_page=100")
	var sizes DropletSizes
	json.Unmarshal([]byte(jsonString), &sizes)
	for _, s := range sizes.Sizes {
		fmt.Println(display.LeftAligned(s.PriceMonth, 5), display.LeftAligned(s.Slug, 20),
			display.LeftAligned(s.Memory, 10),
			display.LeftAligned(s.Disk, 10),
			display.LeftAligned(s.Vcpus, 10), s.Description)
		//fmt.Println(strings.Join(s.Regions, ","))
	}
}
func GetClient(pat string) (*godo.Client, context.Context) {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)
	ctx := context.TODO()

	return client, ctx
}
func ListDroplets(pat string, ip2wolf map[string]string) {
	client, ctx := GetClient(pat)

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, _ := client.Droplets.List(ctx, opt)

	for _, droplet := range droplets {
		//fmt.Println(droplet)
		if len(droplet.Networks.V4) < 2 {
			continue
		}
		ip := droplet.Networks.V4[1].IPAddress
		wolfName := ip2wolf[ip]
		if droplet.Name == "many.pw" {
			continue
		}
		display.DisplayServer(wolfName, droplet.ID, "DO", droplet.Name, ip)
		//fmt.Println(droplet.ID)
	}
}

func ListProducerIps(pat string) []string {
	list := []string{}
	client, ctx := GetClient(pat)

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, _ := client.Droplets.List(ctx, opt)

	for _, droplet := range droplets {
		if len(droplet.Networks.V4) < 2 {
			continue
		}
		if droplet.Name != "producer" {
			continue
		}
		ip := droplet.Networks.V4[1].IPAddress
		list = append(list, ip)
	}
	return list
}
