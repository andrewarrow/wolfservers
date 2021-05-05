package digitalocean

import (
	"fmt"
	"os"
	"strconv"

	"github.com/digitalocean/godo"
)

func CreateDroplet(pat, dropletName, size, key string) {

	imageIdString := os.Getenv("DO_IMAGE")

	dci := godo.DropletCreateImage{
		Slug: "ubuntu-18-04-x64",
	}
	if imageIdString != "" {
		imageId, _ := strconv.Atoi(imageIdString)
		dci.Slug = ""
		dci.ID = imageId
	}

	sshKey := godo.DropletCreateSSHKey{0, key}
	sshList := []godo.DropletCreateSSHKey{sshKey}
	createRequest := &godo.DropletCreateRequest{
		Name:    dropletName,
		Region:  "SFO3",
		Size:    size,
		SSHKeys: sshList,
		Image:   dci,
	}
	client, ctx := GetClient(pat)
	r, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("err: %v\n\n", err)
	}
	fmt.Println(r)
}
