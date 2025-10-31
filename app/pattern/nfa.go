package pattern

import (
	"fmt"
)

// State represents a state in the NFA
type State struct {
	ID        int
	Transitions []Transition
	IsAccepting bool
}

// Transition represents a transition from one state to another
type Transition struct {
	Char     *rune    // nil means epsilon transition
	CharSet  map[rune]bool // for character classes like [a-z] or \d
	IsNegative bool   // if true, CharSet contains characters NOT to match
	Next     *State
}

// NFA represents a Non-deterministic Finite Automaton
type NFA struct {
	Start     *State
	End       *State
	StateCount int
	AllStates map[int]*State
}

// NewState creates a new state
func NewState(isAccepting bool) *State {
	return &State{
		Transitions: make([]Transition, 0),
		IsAccepting: isAccepting,
	}
}

// NewNFA creates a new NFA with start and end states
func NewNFA() *NFA {
	start := NewState(false)
	end := NewState(true)
	return &NFA{
		Start:     start,
		End:       end,
		AllStates: make(map[int]*State),
	}
}

// AddTransition adds a transition between two states
func (s *State) AddTransition(char *rune, charSet map[rune]bool, isNegative bool, next *State) {
	s.Transitions = append(s.Transitions, Transition{
		Char:    char,
		CharSet: charSet,
		IsNegative: isNegative,
		Next:    next,
	})
}

// AddEpsilonTransition adds an epsilon transition
func (s *State) AddEpsilonTransition(next *State) {
	s.AddTransition(nil, nil, false, next)
}

// AddCharTransition adds a character transition
func (s *State) AddCharTransition(char rune, next *State) {
	s.AddTransition(&char, nil, false, next)
}

// AddCharSetTransition adds a character set transition
func (s *State) AddCharSetTransition(charSet map[rune]bool, next *State) {
	s.AddTransition(nil, charSet, false, next)
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

// EpsilonClosure computes epsilon-closure of a set of states
func EpsilonClosure(states []*State) []*State {
	closure := make(map[*State]bool)
	stack := make([]*State, 0)
	
	for _, s := range states {
		closure[s] = true
		stack = append(stack, s)
	}
	
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		
		for _, trans := range current.Transitions {
			if trans.Char == nil && trans.CharSet == nil && trans.Next != nil {
				if !closure[trans.Next] {
					closure[trans.Next] = true
					stack = append(stack, trans.Next)
				}
			}
		}
	}
	
	result := make([]*State, 0, len(closure))
	for s := range closure {
		result = append(result, s)
	}
	
	return result
}

// Matches determines if the input matches the NFA pattern (substring matching)
func (nfa *NFA) Matches(input string) bool {
	runes := []rune(input)
	
	// Try matching at every position in the input string
	for start := 0; start <= len(runes); start++ {
		if matchesFrom(nfa, runes, start) {
			return true
		}
	}
	
	return false
}

// matchesFrom checks if the input matches from a specific starting position
func matchesFrom(nfa *NFA, runes []rune, start int) bool {
	currentStates := EpsilonClosure([]*State{nfa.Start})
	
	for i := start; i < len(runes); i++ {
		char := runes[i]
		nextStates := make(map[*State]bool)
		
		// Find all states reachable with this character
		for _, state := range currentStates {
			for _, trans := range state.Transitions {
				if trans.MatchChar(char) {
					nextStates[trans.Next] = true
				}
			}
		}
		
		// Convert to slice and compute epsilon closure
		stateSlice := make([]*State, 0, len(nextStates))
		for s := range nextStates {
			stateSlice = append(stateSlice, s)
		}
		
		currentStates = EpsilonClosure(stateSlice)
		
		// If we have no states left, this starting position won't work
		if len(currentStates) == 0 {
			return false
		}
		
		// Check if we can accept at this position
		for _, state := range currentStates {
			if state.IsAccepting {
				return true
			}
		}
	}
	
	// Check if any current state is accepting at the end
	for _, state := range currentStates {
		if state.IsAccepting {
			return true
		}
	}
	
	return false
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
			
			intermediate := NewState(false)
			
			switch runes[i+1] {
			case 'd':
				// Digit character class [0-9]
				charSet := make(map[rune]bool)
				for j := '0'; j <= '9'; j++ {
					charSet[j] = true
				}
				current.AddCharSetTransition(charSet, intermediate)
			case 'w':
				// Word character class [a-zA-Z0-9_]
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
			default:
				return nil, fmt.Errorf("unsupported escape sequence: \\%c", runes[i+1])
			}
			
			current = intermediate
			intermediate.AddEpsilonTransition(end)
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
			intermediate.AddEpsilonTransition(end)
			
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
				intermediate := NewState(false)
				current.AddCharTransition(runes[i], intermediate)
				intermediate.AddCharTransition(runes[i], intermediate) // loop back
				intermediate.AddEpsilonTransition(end)
				current = intermediate
				i += 2
			} else if i+1 < len(runes) && runes[i+1] == '*' {
				// Star operator: match zero or more
				intermediate := NewState(false)
				intermediate2 := NewState(false)
				current.AddEpsilonTransition(intermediate)
				current.AddEpsilonTransition(end) // can skip
				intermediate.AddCharTransition(runes[i], intermediate2)
				intermediate2.AddCharTransition(runes[i], intermediate2) // loop back
				intermediate2.AddEpsilonTransition(end)
				current = intermediate2
				i += 2
			} else {
				// Just a regular character
				intermediate := NewState(false)
				current.AddCharTransition(runes[i], intermediate)
				current = intermediate
				intermediate.AddEpsilonTransition(end)
				i++
			}
		}
	}
	
	return nfa, nil
}

