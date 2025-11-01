package pattern

import (
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{
			pattern: "hello",
			input:   "hello world",
			want:    true,
		},
		{
			pattern: `\d apple`,
			input:   "1 apple",
			want:    true,
		},
		{
			pattern: `\d apple`,
			input:   "2 orange",
			want:    false,
		},
		{
			pattern: `\d\d\d apples`,
			input:   "I got 100 apples from the store",
			want:    true,
		},
		{
			pattern: `\d\d\d apples`,
			input:   "I got 1 apple from the store",
			want:    false,
		},
		{
			pattern: `\d \w\w\ws`,
			input:   "4 cats",
			want:    true,
		},
		{
			pattern: `\d \w\w\ws`,
			input:   "1 dog",
			want:    false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {

			nfa, err := ParsePattern(test.pattern)
			if err != nil {
				t.Errorf("ParsePattern(%q) = %v", test.pattern, err)
			}

			got := nfa.Matches(test.input)
			if got != test.want {
				t.Errorf("nfa.Matches(%q, %q) = %v, want %v", test.input, test.pattern, got, test.want)
			}
		})
	}
}
