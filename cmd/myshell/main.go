package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/internal/completer"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
	"github.com/codecrafters-io/shell-starter-go/internal/redirect"
)

// builtinNames lists all shell builtin command names.
var builtinNames = []string{"echo", "exit", "type", "pwd", "cd"}

func main() {
	comp := &completer.Completer{
		Builtins: builtinNames,
	}

	config := &readline.Config{
		Prompt:       "$ ",
		AutoComplete: comp,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			if err == io.EOF || err == readline.ErrInterrupt {
				os.Exit(0)
			}
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		inputParts := parser.ParseInput(input)
		redir := redirect.Parse(inputParts)

		if len(redir.CommandParts) == 0 {
			continue
		}

		command := strings.ToLower(redir.CommandParts[0])

		switch command {
		case "exit":
			handleExit(redir.CommandParts)
		case "echo":
			handleEcho(redir)
		case "type":
			handleType(redir.CommandParts)
		case "pwd":
			handlePwd()
		case "cd":
			handleCd(redir.CommandParts)
		default:
			executeExternal(redir)
		}
	}
}

// handleExit exits the shell. Supports "exit" and "exit 0".
func handleExit(parts []string) {
	if len(parts) == 1 || (len(parts) > 1 && parts[1] == "0") {
		os.Exit(0)
	}
}

// handleEcho writes the arguments to stdout or to a redirected output file.
// Also creates an empty error file if stderr redirection is specified.
func handleEcho(redir redirect.Redirect) {
	if len(redir.CommandParts) < 2 {
		fmt.Println("echo: missing argument")
		return
	}

	output := strings.Join(redir.CommandParts[1:], " ")

	if redir.HasOutput() {
		file, err := redir.OpenOutputFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		fmt.Fprintln(file, output)
		file.Close()
	} else {
		fmt.Println(output)
	}

	// Create empty error file for echo (codecrafters requirement)
	if redir.HasError() {
		file, err := os.Create(redir.ErrorFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		file.Close()
	}
}

// handleType reports whether a command is a builtin or an external executable.
func handleType(parts []string) {
	if len(parts) < 2 {
		return
	}

	target := parts[1]

	// Check if it's a builtin
	for _, name := range builtinNames {
		if target == name {
			fmt.Printf("%s is a shell builtin\n", target)
			return
		}
	}

	// Search PATH for the executable
	pathEnv := os.Getenv("PATH")
	if pathEnv != "" {
		for _, dir := range strings.Split(pathEnv, ":") {
			execPath := filepath.Join(dir, target)
			info, err := os.Stat(execPath)
			if err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
				fmt.Printf("%s is %s\n", target, execPath)
				return
			}
		}
	}

	fmt.Printf("%s: not found\n", target)
}

// handlePwd prints the current working directory.
func handlePwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wd)
}

// handleCd changes the current working directory.
// Supports absolute paths, relative paths, ~ and ~/... for home directory.
func handleCd(parts []string) {
	if len(parts) < 2 {
		fmt.Println("cd: missing argument")
		return
	}

	dir := parts[1]

	// Expand home directory
	if dir == "~" {
		dir = os.Getenv("HOME")
	} else if strings.HasPrefix(dir, "~/") {
		dir = filepath.Join(os.Getenv("HOME"), dir[2:])
	}

	if err := os.Chdir(dir); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
	}
}

// executeExternal runs an external command found in PATH, with I/O redirection support.
func executeExternal(redir redirect.Redirect) {
	commandName := redir.CommandParts[0]

	executable, err := exec.LookPath(commandName)
	if err != nil {
		fmt.Printf("%s: command not found\n", commandName)
		return
	}

	cmd := exec.Command(executable, redir.CommandParts[1:]...)
	cmd.Args[0] = commandName // Use original name, not full path
	cmd.Stdin = os.Stdin

	// Setup stdout
	if redir.HasOutput() {
		file, err := redir.OpenOutputFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		defer file.Close()
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	// Setup stderr
	if redir.HasError() {
		file, err := redir.OpenErrorFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		defer file.Close()
		cmd.Stderr = file
	} else {
		cmd.Stderr = os.Stderr
	}

	cmd.Run()
}
