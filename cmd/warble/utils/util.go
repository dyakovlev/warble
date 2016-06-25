package utils

import (
	"strings"
)

// changes youremailaddress@host.com to y***s@host.com
func ObfuscateEmail(email string) string {
	address, host := strings.Split(email, "@")
	return address[0] + "***" + address[len(address)-1] + "@" + host
}
