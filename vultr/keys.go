package vultr

import (
	"fmt"

	"github.com/vultr/govultr/v2"
)

func ListKeys(pat string) []string {

	list := []string{}
	client, ctx := VultrClient(pat)
	listOptions := &govultr.ListOptions{PerPage: 100}
	i, _, _ := client.SSHKey.List(ctx, listOptions)
	for _, v := range i {
		fmt.Println(v)
		list = append(list, v.ID)
	}
	return list
}
func DeleteKey(pat, id string) {
	client, ctx := VultrClient(pat)
	e := client.SSHKey.Delete(ctx, id)
	fmt.Println(e)
}
func CreateKey(pat, name, pubKey string) {
	req := govultr.SSHKeyReq{}
	req.Name = name
	req.SSHKey = pubKey
	client, ctx := VultrClient(pat)
	e, f := client.SSHKey.Create(ctx, &req)
	fmt.Println(e, f)
}
