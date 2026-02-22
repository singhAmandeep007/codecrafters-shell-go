package parser

// ParseInput parses a shell input string into tokenized arguments,
// handling single quotes, double quotes, and backslash escaping.
func ParseInput(s string) []string {
	var inSingleQuote bool
	var inDoubleQuote bool
	var hasBackslash bool
	var arg string
	var result []string

	for _, char := range s {
		switch char {
		case '\'':
			if hasBackslash && inDoubleQuote {
				arg += "\\"
			}
			if hasBackslash || inDoubleQuote {
				arg += string(char)
			} else {
				inSingleQuote = !inSingleQuote
			}
			hasBackslash = false
		case '"':
			if hasBackslash || inSingleQuote {
				arg += string(char)
			} else {
				inDoubleQuote = !inDoubleQuote
			}
			hasBackslash = false
		case '\\':
			if hasBackslash || inSingleQuote {
				arg += string(char)
				hasBackslash = false
			} else {
				hasBackslash = true
			}
		case ' ':
			if hasBackslash && inDoubleQuote {
				arg += "\\"
			}
			if hasBackslash || inSingleQuote || inDoubleQuote {
				arg += string(char)
			} else if arg != "" {
				result = append(result, arg)
				arg = ""
			}
			hasBackslash = false
		default:
			if inDoubleQuote && hasBackslash {
				arg += "\\"
			}
			arg += string(char)
			hasBackslash = false
		}
	}

	if arg != "" {
		result = append(result, arg)
	}

	return result
}
