package completer

import "testing"

func TestFindLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{name: "empty slice", input: nil, expected: ""},
		{name: "single string", input: []string{"abc"}, expected: "abc"},
		{name: "common prefix", input: []string{"abc", "abd"}, expected: "ab"},
		{name: "no common prefix", input: []string{"abc", "xyz"}, expected: ""},
		{name: "all identical", input: []string{"echo", "echo", "echo"}, expected: "echo"},
		{name: "one char prefix", input: []string{"exit", "echo", "env"}, expected: "e"},
		{name: "longer prefix", input: []string{"export", "expr", "expand"}, expected: "exp"},
		{name: "empty in slice", input: []string{"abc", ""}, expected: ""},
		{name: "prefix is shortest", input: []string{"ab", "abc", "abcd"}, expected: "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindLongestCommonPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("FindLongestCommonPrefix(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindMatches_Builtins(t *testing.T) {
	c := &Completer{
		Builtins: []string{"echo", "exit", "type", "pwd", "cd"},
	}

	tests := []struct {
		name     string
		prefix   string
		expected []string
	}{
		{name: "matches e prefix", prefix: "e", expected: []string{"echo", "exit"}},
		{name: "matches echo exactly", prefix: "echo", expected: []string{"echo"}},
		{name: "matches type", prefix: "ty", expected: []string{"type"}},
		{name: "no match", prefix: "zzz", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.FindMatches(tt.prefix)
			builtinSet := map[string]bool{
				"echo": true, "exit": true, "type": true, "pwd": true, "cd": true,
			}
			var builtinResults []string
			for _, m := range result {
				if builtinSet[m] {
					builtinResults = append(builtinResults, m)
				}
			}
			if len(builtinResults) != len(tt.expected) {
				t.Errorf("FindMatches(%q) builtins = %v, want %v", tt.prefix, builtinResults, tt.expected)
				return
			}
			for i, v := range builtinResults {
				if v != tt.expected[i] {
					t.Errorf("FindMatches(%q) builtins = %v, want %v", tt.prefix, builtinResults, tt.expected)
					return
				}
			}
		})
	}
}
