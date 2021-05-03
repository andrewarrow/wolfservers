package vultr

import "fmt"

func RemoveServer(id string) {
	client, ctx := VultrClient()
	e := client.Instance.Delete(ctx, id)
	fmt.Println(e)
}
