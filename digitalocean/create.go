package digitalocean

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"

	"github.com/digitalocean/godo"
)

func CreateDroplet(size, key string) {
	b := make([]byte, 16)
	rand.Read(b)
	dropletName := fmt.Sprintf("wolf-%X", b[4:6])

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
	client, ctx := GetClient()
	r, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("err: %v\n\n", err)
	}
	fmt.Println(r)
}
