package keys

import (
	"crypto/rand"
	"fmt"
)

func WolfName(prefix string) string {
	b := make([]byte, 16)
	rand.Read(b)
	name := fmt.Sprintf(prefix+"-%X", b[4:6])
	return name
}
