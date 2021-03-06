package linode

import (
	"fmt"
	"os"

	"github.com/andrewarrow/wolfservers/keys"
	"github.com/linode/linodego"
)

func CreateServer(pat, name string) {
	client, ctx := LinodeClient(pat)
	instanceOptions := linodego.InstanceCreateOptions{
		Label:           name,
		RootPass:        keys.RootPass(),
		AuthorizedUsers: []string{os.Getenv("LINODE_USER")},
		Image:           "linode/ubuntu18.04",
		//Type:            "g6-standard-4",
		Type: "g6-standard-2",
		//Type:   "g6-standard-1",
		Region: "us-west",
	}

	res, err := client.CreateInstance(ctx, instanceOptions)

	fmt.Println(res, err)
}

func CreateSshKey(pat, name, pubKey string) {
	client, ctx := LinodeClient(pat)
	options := linodego.SSHKeyCreateOptions{
		Label:  name,
		SSHKey: pubKey,
	}

	res, err := client.CreateSSHKey(ctx, options)

	fmt.Println(res, err)
}
