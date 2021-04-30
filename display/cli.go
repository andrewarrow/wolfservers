package display

import (
	"fmt"
	"strings"
)

func DisplayServer(name, id, provider, label, ip interface{}) {
	fmt.Printf("%v %s %v %s %s\n",
		name,
		LeftAligned(id, 10),
		LeftAligned(provider, 7),
		LeftAligned(label, 10),
		LeftAligned(ip, 30))
}

func LeftAligned(thing interface{}, size int) string {
	s := fmt.Sprintf("%v", thing)

	if len(s) > size {
		return s[0:size]
	}
	fill := size - len(s)
	spaces := []string{}
	for {
		spaces = append(spaces, " ")
		if len(spaces) >= fill {
			break
		}
	}
	return s + strings.Join(spaces, "")
}
