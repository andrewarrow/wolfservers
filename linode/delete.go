package linode

import (
	"fmt"
)

func DeleteSshKey(pat string, id int) {
	client, ctx := LinodeClient(pat)
	err := client.DeleteSSHKey(ctx, id)
	fmt.Println(err)
}
