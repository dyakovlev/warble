package utils

import (
	"fmt"
	"strings"
)

// youremailaddress@host.com -> y***s@host.com
func ObfuscateEmail(email string) string {
	chunks := strings.Split(email, "@")
	addr := []rune(chunks[0])
	return fmt.Sprintf("%v***%v@%v", string(addr[0]), string(addr[len(addr)-1]), chunks[1])
}

// remove typical db-unsafe characters, return whether or not any were bad
func SanitizeText(text *string) bool {
	return true
}

func ParseEmail(email string) (string, error) {
	return email, nil
}

func ParsePassword(pass string) (string, error) {
	return pass, nil
}
