package pattern

import (
	"strings"
)

const WordCharacterClass = `\w`

func ContainsWordCharacterClass(pattern string) bool {
	return strings.Contains(pattern, WordCharacterClass)
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
