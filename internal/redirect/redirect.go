package redirect

import "os"

// Redirect holds parsed I/O redirection information for a shell command.
type Redirect struct {
	OutputFile   string   // File path for stdout redirection (>, 1>, >>, 1>>)
	ErrorFile    string   // File path for stderr redirection (2>, 2>>)
	AppendOutput bool     // True if stdout should append (>>, 1>>)
	AppendError  bool     // True if stderr should append (2>>)
	CommandParts []string // The command and arguments without redirect operators
}

// Parse scans inputParts for the first I/O redirection operator and returns
// a Redirect with the parsed information. Only one redirect operator per
// command is supported (consistent with codecrafters requirements).
func Parse(inputParts []string) Redirect {
	var r Redirect

	for i, part := range inputParts {
		switch part {
		case ">", "1>":
			if i+1 < len(inputParts) {
				r.OutputFile = inputParts[i+1]
				r.CommandParts = inputParts[:i]
				r.AppendOutput = false
				return r
			}
		case ">>", "1>>":
			if i+1 < len(inputParts) {
				r.OutputFile = inputParts[i+1]
				r.CommandParts = inputParts[:i]
				r.AppendOutput = true
				return r
			}
		case "2>":
			if i+1 < len(inputParts) {
				r.ErrorFile = inputParts[i+1]
				r.CommandParts = inputParts[:i]
				r.AppendError = false
				return r
			}
		case "2>>":
			if i+1 < len(inputParts) {
				r.ErrorFile = inputParts[i+1]
				r.CommandParts = inputParts[:i]
				r.AppendError = true
				return r
			}
		}
	}

	// No redirect found
	r.CommandParts = inputParts
	return r
}

// HasOutput returns true if stdout is being redirected to a file.
func (r *Redirect) HasOutput() bool {
	return r.OutputFile != ""
}

// HasError returns true if stderr is being redirected to a file.
func (r *Redirect) HasError() bool {
	return r.ErrorFile != ""
}

// OpenOutputFile opens the output redirect file with appropriate flags.
// Returns nil if no output redirection is configured.
func (r *Redirect) OpenOutputFile() (*os.File, error) {
	if r.OutputFile == "" {
		return nil, nil
	}
	if r.AppendOutput {
		return os.OpenFile(r.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(r.OutputFile)
}

// OpenErrorFile opens the error redirect file with appropriate flags.
// Returns nil if no error redirection is configured.
func (r *Redirect) OpenErrorFile() (*os.File, error) {
	if r.ErrorFile == "" {
		return nil, nil
	}
	if r.AppendError {
		return os.OpenFile(r.ErrorFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(r.ErrorFile)
}
