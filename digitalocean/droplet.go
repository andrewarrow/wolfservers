package digitalocean

type DropletSize struct {
	Slug       string
	Memory     string
	Disk       string
	PriceMonth int
	Regions    []string
	Available  bool
}

type DropletSizes struct {
	Sizes []DropletSize
}
