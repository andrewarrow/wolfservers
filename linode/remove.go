package linode

import (
	"fmt"
)

func RemoveServer(pat string, id int) {
	client, ctx := LinodeClient(pat)

	err := client.DeleteInstance(ctx, id)

	fmt.Println(err)
}
