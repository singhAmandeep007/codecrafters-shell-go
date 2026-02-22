package completer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Completer implements readline.AutoCompleter for shell tab completion.
// It completes builtin commands and PATH executables, supports longest
// common prefix (LCP) completion, and displays all matches on double-TAB.
type Completer struct {
	Builtins    []string
	lastLine    string
	lastMatches []string
	tabCount    int
}

// Do is called by readline on each TAB press. It returns completion
// candidates and the length of the text to replace.
func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])

	// Track consecutive TAB presses on the same input
	if lineStr == c.lastLine {
		c.tabCount++
	} else {
		c.tabCount = 1
		c.lastLine = lineStr
	}

	matches := c.FindMatches(lineStr)
	c.lastMatches = matches

	// No matches: ring the bell
	if len(matches) == 0 {
		fmt.Print("\x07")
		c.tabCount = 0
		return nil, len(lineStr)
	}

	// Single match: complete with trailing space
	if len(matches) == 1 {
		c.tabCount = 0
		completion := matches[0][len(lineStr):] + " "
		return [][]rune{[]rune(completion)}, len(lineStr)
	}

	// Multiple matches: try LCP completion
	lcp := FindLongestCommonPrefix(matches)
	if len(lcp) > len(lineStr) {
		c.tabCount = 0
		completion := lcp[len(lineStr):]
		return [][]rune{[]rune(completion)}, len(lineStr)
	}

	// LCP equals current input - cannot complete further
	if c.tabCount == 1 {
		// First TAB: ring the bell
		fmt.Print("\x07")
		return nil, len(lineStr)
	} else if c.tabCount >= 2 {
		// Second TAB: display all matches sorted
		sort.Strings(matches)
		os.Stdout.WriteString("\n" + strings.Join(matches, "  ") + "\n$ " + lineStr)
		os.Stdout.Sync()
		c.tabCount = 0
		return nil, len(lineStr)
	}

	return nil, len(lineStr)
}

// FindMatches collects matching builtins and PATH executables for the given prefix.
func (c *Completer) FindMatches(prefix string) []string {
	seen := make(map[string]bool)
	var matches []string

	// Match builtins
	for _, builtin := range c.Builtins {
		if strings.HasPrefix(builtin, prefix) {
			matches = append(matches, builtin)
			seen[builtin] = true
		}
	}

	// Match PATH executables
	pathEnv := os.Getenv("PATH")
	if pathEnv != "" {
		for _, dir := range strings.Split(pathEnv, ":") {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				name := entry.Name()
				if !strings.HasPrefix(name, prefix) || seen[name] {
					continue
				}
				fullPath := filepath.Join(dir, name)
				info, err := os.Stat(fullPath)
				if err != nil {
					continue
				}
				if !info.IsDir() && info.Mode()&0111 != 0 {
					matches = append(matches, name)
					seen[name] = true
				}
			}
		}
	}

	return matches
}

// FindLongestCommonPrefix returns the longest common prefix of all strings.
func FindLongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if len(prefix) == 0 {
				return ""
			}
		}
	}
	return prefix
}
