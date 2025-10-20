package pattern

import (
	"strings"
)

const DigitCharacterClass = `\d`

func ContainsDigitCharacterClass(pattern string) bool {
	return strings.Contains(pattern, DigitCharacterClass)
}

func ContainsDigit(input []byte) bool {
	for _, r := range string(input) {
		if IsDigit(r) {
			return true
		}
	}
	return false
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
