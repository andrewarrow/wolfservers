package vultr

import (
	"github.com/vultr/govultr/v2"
)

func ListKeys() []string {

	list := []string{}
	client, ctx := VultrClient()
	listOptions := &govultr.ListOptions{PerPage: 100}
	i, _, _ := client.SSHKey.List(ctx, listOptions)
	for _, v := range i {
		list = append(list, v.ID)
	}
	return list
}
func DeleteKey(id string) {
	client, ctx := VultrClient()
	client.SSHKey.Delete(ctx, id)
}
func CreateKey(name, pubKey string) {
	req := govultr.SSHKeyReq{}
	req.Name = name
	req.SSHKey = pubKey
	client, ctx := VultrClient()
	client.SSHKey.Create(ctx, &req)
}
