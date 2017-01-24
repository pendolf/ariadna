package common

import (
	"fmt"
	"regexp"
	"strings"
)

var Latinre *regexp.Regexp

func init() {
	Latinre, _ = regexp.Compile("^[a-zA-Z]")
}

func NormalizeAddress(address string) string {
	if strings.Contains(address, "улица") {
		return fmt.Sprintf("улица %s", strings.Replace(address, "улица", "", -1))

	}
	if strings.Contains(address, "проспект") {
		return fmt.Sprintf("проспект %s", strings.Replace(address, "проспект", "", -1))

	}
	if strings.Contains(address, "переулок") {
		return fmt.Sprintf("переулок %s", strings.Replace(address, "переулок", "", -1))

	}
	return address
}
func CleanAddress(address string) string {
	if strings.Contains(address, "улица") {
		return strings.Replace(address, "улица", "", -1)

	}
	if strings.Contains(address, "проспект") {
		return strings.Replace(address, "проспект", "", -1)

	}
	if strings.Contains(address, "переулок") {
		return strings.Replace(address, "переулок", "", -1)

	}
	return address
}
