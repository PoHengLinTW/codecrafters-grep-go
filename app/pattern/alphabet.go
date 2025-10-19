package pattern

func IsAlpha(r rune) bool {
	return IsLowerAlpha(r) || IsUpperAlpha(r)
}

func IsLowerAlpha(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func IsUpperAlpha(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
