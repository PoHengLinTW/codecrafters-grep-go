package pattern

import "strings"

func IsPositiveCharacterGroup(pattern string) bool {
	return strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
}

func ContainsPositiveCharacterGroup(input []byte, pattern string) bool {
	if !IsPositiveCharacterGroup(pattern) {
		return false
	}

	groupMap := make(map[rune]bool)

	for _, r := range pattern[1 : len(pattern)-1] {
		groupMap[r] = true
	}

	for _, r := range string(input) {
		_, ok := groupMap[r]
		if ok {
			return true
		}
	}

	return false
}
