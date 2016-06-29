package utils

import (
	"fmt"
	"strings"
)

// changes youremailaddress@host.com to y***s@host.com
func ObfuscateEmail(email string) string {
	chunks := strings.Split(email, "@")
	addr := []rune(chunks[0])
	return fmt.Sprintf("%v***%v@%v", string(addr[0]), string(addr[len(addr)-1]), chunks[1])
}
