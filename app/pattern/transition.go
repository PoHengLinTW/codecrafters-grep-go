package pattern

// Transition represents a transition from one state to another
type Transition struct {
	Char     *rune    // nil means epsilon transition
	CharSet  map[rune]bool // for character classes like [a-z] or \d
	IsNegative bool   // if true, CharSet contains characters NOT to match
	Next     *State
}

// MatchChar checks if a character matches this transition
func (t *Transition) MatchChar(char rune) bool {
	if t.Char != nil {
		return *t.Char == char
	}
	if t.CharSet != nil {
		if t.IsNegative {
			return !t.CharSet[char]
		}
		return t.CharSet[char]
	}
	return false
}
