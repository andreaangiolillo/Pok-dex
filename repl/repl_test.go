package repl

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		out := CleanInput(c.input)
		if len(out) != len(c.expected) {
			t.Errorf("got: %d, want: %d", len(out), len(c.expected))
		}

		for i:= range out{
			if out[i] != c.expected[i] {
				t.Errorf("got: %s, want: %s", out[i], c.expected[i])
				}
			}
	}
}

