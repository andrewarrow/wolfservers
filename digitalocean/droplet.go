package digitalocean

type DropletSize struct {
	Slug        string
	Memory      string
	Disk        string
	PriceMonth  float64 `json:"price_monthly"`
	Regions     []string
	Available   bool
	Description string
}

type DropletSizes struct {
	Sizes []DropletSize
}
