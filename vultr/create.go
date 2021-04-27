package vultr

import (
	"fmt"

	"github.com/vultr/govultr/v2"
)

func BoolToBoolPtr(value bool) *bool {
	b := value
	return &b
}

func CreateServer() {
	client, ctx := VultrClient()
	instanceOptions := &govultr.InstanceCreateReq{
		Label:      "test",
		Hostname:   "test",
		Backups:    "false",
		EnableIPv6: BoolToBoolPtr(false),
		OsID:       362,
		Plan:       "vc2-1c-2gb",
		Region:     "ewr",
	}

	res, err := client.Instance.Create(ctx, instanceOptions)

	fmt.Println(res, err)
}
