package linode

import (
	"fmt"
)

func DeleteSshKey(id int) {
	client, ctx := LinodeClient()
	err := client.DeleteSSHKey(ctx, id)
	fmt.Println(err)
}
