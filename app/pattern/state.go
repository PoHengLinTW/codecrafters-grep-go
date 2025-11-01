package pattern

// State represents a state in the NFA
type State struct {
	ID        int
	Transitions []Transition
	IsAccepting bool
}

// NewState creates a new state
func NewState(isAccepting bool) *State {
	return &State{
		Transitions: make([]Transition, 0),
		IsAccepting: isAccepting,
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
