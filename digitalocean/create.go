package digitalocean

import (
	"crypto/rand"
	"fmt"

	"github.com/digitalocean/godo"
)

func CreateDroplet(size, key string) {
	b := make([]byte, 16)
	rand.Read(b)
	dropletName := fmt.Sprintf("wolf-%X", b[4:6])

	sshKey := godo.DropletCreateSSHKey{0, key}
	sshList := []godo.DropletCreateSSHKey{sshKey}
	createRequest := &godo.DropletCreateRequest{
		Name:    dropletName,
		Region:  "SFO3",
		Size:    size,
		SSHKeys: sshList,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-18-04-x64",
		},
	}

	client, ctx := GetClient()
	_, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("err: %v\n\n", err)

	}
}
