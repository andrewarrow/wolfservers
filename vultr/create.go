package vultr

import (
	"fmt"

	"github.com/vultr/govultr/v2"
)

func BoolToBoolPtr(value bool) *bool {
	b := value
	return &b
}

func CreateServer(pat, name string) {
	keys := ListKeys(pat)
	client, ctx := VultrClient(pat)
	instanceOptions := &govultr.InstanceCreateReq{
		Label:      name,
		Hostname:   name,
		Backups:    "disabled",
		EnableIPv6: BoolToBoolPtr(false),
		OsID:       270,
		//Plan:       "vc2-4c-8gb",
		Plan:    "vc2-2c-4gb",
		SSHKeys: keys,
		Region:  "lax",
	}

	res, err := client.Instance.Create(ctx, instanceOptions)

	fmt.Println(res, err)
}
