package digitalocean

import "fmt"

func RemoveDroplet(pat string, id int) {
	client, ctx := GetClient(pat)

	resp, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != 204 {
		fmt.Println(resp.StatusCode)
	}
	fmt.Println(resp)
}
