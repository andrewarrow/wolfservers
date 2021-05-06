package files

import (
	"io/ioutil"
	"strings"
)

func RemoveComments(in string) {
	buff := []string{}
	b, _ := ioutil.ReadFile(in)
	s := string(b)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "//") == false {
			buff = append(buff, line)
		}
	}

	ioutil.WriteFile(in, []byte(strings.Join(buff, "\n")), 0755)

}
