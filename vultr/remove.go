package vultr

import "fmt"

func RemoveServer(pat, id string) {
	client, ctx := VultrClient(pat)
	e := client.Instance.Delete(ctx, id)
	fmt.Println(e)
}
