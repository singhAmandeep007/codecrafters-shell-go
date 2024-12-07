package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	// This package provides functions to manipulate filename paths in a way that's compatible with the operating system where the program is running.
	// It's particularly useful for building file paths using slashes on Unix-like systems or backslashes on Windows.
	// In the context of the shell we're building, this package will be used to construct file paths when searching for executable files in the directories specified by the PATH environment variable.
	"path/filepath"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		// reads a line of input from the user and stores it in the variable input.
		// The underscore _ could be used to ignore any error returned by ReadString.
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			os.Exit(1)
		}

		// The strings.TrimSpace function is used to remove any leading or trailing whitespace from the input string.
		input = strings.TrimSpace(input)
		// The strings.Split function is used to split the input string into a slice of strings.
		inputParts := strings.Split(input, " ")

		// If the user enters the exit command, the shell will exit.
		if input == "exit 0" {
			os.Exit(0)
		}

		if inputParts[0] == "echo" {
			// The first word is the command name, and the rest of the words are the arguments.
			// The arguments are joined together with a space character and printed to the console.
			fmt.Printf("%s\n", strings.Join(inputParts[1:], " "))
			continue
		}

		if inputParts[0] == "type" {
			switch inputParts[1] {
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			case "exit":
				fmt.Println("exit is a shell builtin")
			default:
				command := inputParts[1]
				// The os.Getenv function is used to retrieve the value of the PATH environment variable.
				pathEnv := os.Getenv("PATH")
				// The isFound variable is used to keep track of whether the command was found in any of the directories.
				isFound := false

				if pathEnv != "" {
					// Split PATH into directories
					envPaths := strings.Split(pathEnv, ":")
					// Search for the command in each directory
					for _, dir := range envPaths {
						// The filepath.Join function is used to construct the full path to the executable file by joining the directory path and the command name.
						exec := filepath.Join(dir, command)
						// The os.Stat function is used to check if the file exists.
						// If the file exists, the command is printed along with the full path to the executable file.
						if _, err := os.Stat(exec); err == nil {
							fmt.Printf("%v is %v\n", command, exec)
							isFound = true
							break
						}
					}
					if !isFound {
						fmt.Printf("%s: not found\n", command)
					}
				} else {
					fmt.Printf("%s: not found\n", command)
				}
			}
			continue
		}

		// DEFAULT
		// prints a message indicating that the command is not found.
		// The input[:len(input)-1] part removes the newline character from the end of the input string, ensuring the command name is printed correctly without an extra line break.
		// This change allows the shell to handle invalid commands by displaying a message in the format <command_name>: command not found
		fmt.Printf("%s: command not found\n", input[:])
	}
}
