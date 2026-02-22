package redirect

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected Redirect
	}{
		{
			name:  "no redirect",
			input: []string{"echo", "hello"},
			expected: Redirect{
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stdout redirect >",
			input: []string{"echo", "hello", ">", "file.txt"},
			expected: Redirect{
				OutputFile:   "file.txt",
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stdout redirect 1>",
			input: []string{"echo", "hello", "1>", "file.txt"},
			expected: Redirect{
				OutputFile:   "file.txt",
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stdout append >>",
			input: []string{"echo", "hello", ">>", "file.txt"},
			expected: Redirect{
				OutputFile:   "file.txt",
				AppendOutput: true,
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stdout append 1>>",
			input: []string{"echo", "hello", "1>>", "file.txt"},
			expected: Redirect{
				OutputFile:   "file.txt",
				AppendOutput: true,
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stderr redirect 2>",
			input: []string{"echo", "hello", "2>", "err.txt"},
			expected: Redirect{
				ErrorFile:    "err.txt",
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "stderr append 2>>",
			input: []string{"echo", "hello", "2>>", "err.txt"},
			expected: Redirect{
				ErrorFile:    "err.txt",
				AppendError:  true,
				CommandParts: []string{"echo", "hello"},
			},
		},
		{
			name:  "redirect with no args before operator",
			input: []string{">", "file.txt"},
			expected: Redirect{
				OutputFile:   "file.txt",
				CommandParts: []string{},
			},
		},
		{
			name:  "single element no redirect",
			input: []string{"pwd"},
			expected: Redirect{
				CommandParts: []string{"pwd"},
			},
		},
		{
			name:  "empty input",
			input: []string{},
			expected: Redirect{
				CommandParts: []string{},
			},
		},
		{
			name:  "redirect operator without target file",
			input: []string{"echo", "hello", ">"},
			expected: Redirect{
				CommandParts: []string{"echo", "hello", ">"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Parse(%v)\n  got:  %+v\n  want: %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRedirectHasOutput(t *testing.T) {
	r := Redirect{OutputFile: "out.txt"}
	if !r.HasOutput() {
		t.Error("HasOutput() should return true when OutputFile is set")
	}

	r2 := Redirect{}
	if r2.HasOutput() {
		t.Error("HasOutput() should return false when OutputFile is empty")
	}
}

func TestRedirectHasError(t *testing.T) {
	r := Redirect{ErrorFile: "err.txt"}
	if !r.HasError() {
		t.Error("HasError() should return true when ErrorFile is set")
	}

	r2 := Redirect{}
	if r2.HasError() {
		t.Error("HasError() should return false when ErrorFile is empty")
	}
}
