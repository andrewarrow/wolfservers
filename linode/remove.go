package linode

import (
	"fmt"
)

func RemoveServer(id int) {
	client, ctx := LinodeClient()

	err := client.DeleteInstance(ctx, id)

	fmt.Println(err)
}
