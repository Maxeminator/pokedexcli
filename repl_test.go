package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Gergol Rogan",
			expected: []string{"gergol", "rogan"},
		},
		{
			input:    "ASDQWEzxc zxcasdQWE",
			expected: []string{"asdqwezxc", "zxcasdqwe"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length mismatched expected %v, actual %v", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Word mismatch expected %v, actual %v", expectedWord, actual)
			}

		}
	}
}
