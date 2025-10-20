package pattern

import (
	"strings"
)

const WordCharacterClass = `\w`

func ContainsWordCharacterClass(input string) bool {
	return strings.Contains(input, WordCharacterClass)
}

func ContainsWord(input []byte) bool {
	for _, r := range string(input) {
		if IsAlphaNumeric(r) || r == '_' {
			return true
		}
	}
	return false
}

func IsAlphaNumeric(r rune) bool {
	return IsDigit(r) || IsAlpha(r)
}
