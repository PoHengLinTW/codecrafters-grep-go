package pattern

// NFA represents a Non-deterministic Finite Automaton
type NFA struct {
	Start     *State
	End       *State
	StateCount int
	AllStates map[int]*State
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

