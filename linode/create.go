package linode

import (
	"fmt"
	"os"

	"github.com/linode/linodego"
)

func CreateServer(name string) {
	client, ctx := LinodeClient()
	instanceOptions := linodego.InstanceCreateOptions{
		Label:           name,
		AuthorizedUsers: []string{os.Getenv("LINODE_USER")},
		Image:           "linode/ubuntu18.04",
		Type:            "g6-standard-4",
		Region:          "us-west",
	}

	res, err := client.CreateInstance(ctx, instanceOptions)

	fmt.Println(res, err)
}
