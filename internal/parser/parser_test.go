package parser

import (
	"reflect"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple command",
			input:    "echo hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "multiple arguments",
			input:    "echo hello world foo",
			expected: []string{"echo", "hello", "world", "foo"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "single word",
			input:    "pwd",
			expected: []string{"pwd"},
		},
		{
			name:     "extra spaces between words",
			input:    "echo   hello   world",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "leading and trailing spaces",
			input:    "  echo hello  ",
			expected: []string{"echo", "hello"},
		},

		// Single quote tests
		{
			name:     "single quotes preserve spaces",
			input:    "echo 'hello world'",
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "single quotes preserve double quotes",
			input:    `echo 'say "hi"'`,
			expected: []string{"echo", `say "hi"`},
		},
		{
			name:     "single quotes preserve backslash",
			input:    `echo 'path\\to'`,
			expected: []string{"echo", `path\\to`},
		},
		{
			name:     "adjacent single quoted strings",
			input:    "echo 'hello'' world'",
			expected: []string{"echo", "hello world"},
		},

		// Double quote tests
		{
			name:     "double quotes preserve spaces",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "double quotes preserve single quotes",
			input:    `echo "it's fine"`,
			expected: []string{"echo", "it's fine"},
		},
		{
			name:     "escaped double quote inside double quotes",
			input:    `echo "say \"hi\""`,
			expected: []string{"echo", `say "hi"`},
		},
		{
			name:     "escaped backslash inside double quotes",
			input:    `echo "back\\slash"`,
			expected: []string{"echo", `back\slash`},
		},

		// Backslash tests (outside quotes)
		{
			name:     "backslash escapes space",
			input:    `echo hello\ world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "backslash escapes backslash",
			input:    `echo hello\\world`,
			expected: []string{"echo", `hello\world`},
		},
		{
			name:     "backslash escapes single quote",
			input:    `echo it\'s`,
			expected: []string{"echo", "it's"},
		},

		// Redirection-related parsing (parser just tokenizes, doesn't interpret)
		{
			name:     "redirect operator as separate token",
			input:    "echo hello > file.txt",
			expected: []string{"echo", "hello", ">", "file.txt"},
		},
		{
			name:     "stderr redirect",
			input:    "echo hello 2> err.txt",
			expected: []string{"echo", "hello", "2>", "err.txt"},
		},
		{
			name:     "append operator",
			input:    "echo hello >> file.txt",
			expected: []string{"echo", "hello", ">>", "file.txt"},
		},

		// Mixed quotes
		{
			name:     "mixed single and double quotes",
			input:    `echo "hello" 'world'`,
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "quoted and unquoted adjacent",
			input:    `echo hel"lo wo"rld`,
			expected: []string{"echo", "hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseInput(%q)\n  got:  %v\n  want: %v", tt.input, result, tt.expected)
			}
		})
	}
}
