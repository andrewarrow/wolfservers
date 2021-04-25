package digitalocean

import (
	"encoding/json"
	"fmt"

	"github.com/andrewarrow/wolfservers/display"
	"github.com/andrewarrow/wolfservers/network"
)

type Image struct {
	Id   int
	Slug string
	Name string
}
type Images struct {
	Images []Image `json:"images"`
}

func ListImages(page int) {
	jsonString := network.DoGet(fmt.Sprintf("v2/images?page=%d&per_page=200", page))
	var images Images
	json.Unmarshal([]byte(jsonString), &images)
	for _, s := range images.Images {
		fmt.Println(display.LeftAligned(s.Id, 20),
			display.LeftAligned(s.Slug, 30),
			display.LeftAligned(s.Name, 30))
	}
}
