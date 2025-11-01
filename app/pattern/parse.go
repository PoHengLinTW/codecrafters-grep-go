package pattern

import "fmt"

func handleNumberCharacterClass(current *State) *State {
	intermediate := NewState(false)
	charSet := make(map[rune]bool)
	for j := '0'; j <= '9'; j++ {
		charSet[j] = true
	}
	current.AddCharSetTransition(charSet, intermediate)
	return intermediate
}

func handleWordCharacterClass(current *State) *State {
	intermediate := NewState(false)
	charSet := make(map[rune]bool)

	for j := 'a'; j <= 'z'; j++ {
		charSet[j] = true
	}
	for j := 'A'; j <= 'Z'; j++ {
		charSet[j] = true
	}
	for j := '0'; j <= '9'; j++ {
		charSet[j] = true
	}
	charSet['_'] = true

	current.AddCharSetTransition(charSet, intermediate)
	return intermediate
}

func handlePlusOperator(char rune, current *State) *State {
	intermediate := NewState(false)

	current.AddCharTransition(char, intermediate)
	intermediate.AddCharTransition(char, intermediate) // loop back

	return intermediate
}

func handleStarOperator(char rune, current *State, end *State) *State {
	intermediate := NewState(false)
	intermediate2 := NewState(false)

	current.AddEpsilonTransition(intermediate)
	current.AddEpsilonTransition(end) // can skip
	intermediate.AddCharTransition(char, intermediate2)
	intermediate2.AddCharTransition(char, intermediate2) // loop back

	return intermediate2
}

func handleRegularCharacter(char rune, current *State) *State {
	intermediate := NewState(false)
	current.AddCharTransition(char, intermediate)
	return intermediate
}

// ParsePattern builds an NFA from a regex pattern string
func ParsePattern(pattern string) (*NFA, error) {
	nfa := NewNFA()
	current := nfa.Start
	end := nfa.End
	runes := []rune(pattern)
	i := 0
	
	for i < len(runes) {
		switch runes[i] {
		case '\\':
			// Escape sequence
			if i+1 >= len(runes) {
				return nil, fmt.Errorf("unexpected end after backslash")
			}
			
			switch runes[i+1] {
			case 'd':
				current = handleNumberCharacterClass(current)
			case 'w':
				current = handleWordCharacterClass(current)
			default:
				return nil, fmt.Errorf("unsupported escape sequence: \\%c", runes[i+1])
			}
			
			i += 2
			
		case '[':
			// Character group
			i++ // skip '['
			isNegative := false
			
			if i < len(runes) && runes[i] == '^' {
				isNegative = true
				i++
			}
			
			charSet := make(map[rune]bool)
			
			for i < len(runes) && runes[i] != ']' {
				charSet[runes[i]] = true
				i++
			}
			
			if i >= len(runes) {
				return nil, fmt.Errorf("unterminated character group")
			}
			
			i++ // skip ']'
			
			// Create intermediate state
			intermediate := NewState(false)
			
			// Use AddTransition with isNegative flag for proper negative matching
			current.AddTransition(nil, charSet, isNegative, intermediate)
			
			current = intermediate
			
		case '^':
			// Beginning of string anchor - for now, skip it
			i++
			
		case '$':
			// End of string anchor - for now, skip it
			i++
			
		default:
			// Regular character
			if i+1 < len(runes) && runes[i+1] == '+' {
				// Plus operator: match one or more
				current = handlePlusOperator(runes[i], current)
				i += 2
			} else if i+1 < len(runes) && runes[i+1] == '*' {
				// Star operator: match zero or more
				current = handleStarOperator(runes[i], current, end)
				i += 2
			} else {
				// Just a regular character
				current = handleRegularCharacter(runes[i], current)
				i++
			}
		}
	}

	current.AddEpsilonTransition(end)
	
	return nfa, nil
}

