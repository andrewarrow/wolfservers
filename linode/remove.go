package linode

import (
	"fmt"
)

func RemoveServer(id int) {
	client, ctx := LinodeClient()

	fmt.Println(id)
	err := client.DeleteInstance(ctx, id)

	fmt.Println(err)
}
