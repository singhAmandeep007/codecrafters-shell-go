package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		command := inputParts[0]

		// If the user enters the exit command, the shell will exit.
		if input == "exit 0" {
			os.Exit(0)
		}

		if command == "echo" {
			// The first word is the command name, and the rest of the words are the arguments.
			// The arguments are joined together with a space character and printed to the console.
			fmt.Printf("%s\n", strings.Join(inputParts[1:], " "))
			continue
		}

		if command == "type" {
			switch inputParts[1] {
			case "echo":
				fmt.Println("echo is a shell builtin")
			case "type":
				fmt.Println("type is a shell builtin")
			case "exit":
				fmt.Println("exit is a shell builtin")
			case "pwd":
				fmt.Println("pwd is a shell builtin")
			case "cd":
				fmt.Println("cd is a shell builtin")
			default:
				commandArg := inputParts[1]
				// The os.Getenv function is used to retrieve the value of the PATH environment variable.
				pathEnv := os.Getenv("PATH")
				// The isFound variable is used to keep track of whether the commandArg was found in any of the directories.
				isFound := false

				if pathEnv != "" {
					// Split PATH into directories
					envPaths := strings.Split(pathEnv, ":")
					// Search for the commandArg in each directory
					for _, dir := range envPaths {
						// The filepath.Join function is used to construct the full path to the executable file by joining the directory path and the command name.
						exec := filepath.Join(dir, commandArg)
						// The os.Stat function is used to check if the file exists.
						// If the file exists, the commandArg is printed along with the full path to the executable file.
						if _, err := os.Stat(exec); err == nil {
							fmt.Printf("%v is %v\n", commandArg, exec)
							isFound = true
							break
						}
					}
					if !isFound {
						fmt.Printf("%s: not found\n", commandArg)
					}
				} else {
					fmt.Printf("%s: not found\n", commandArg)
				}
			}
			continue
		}

		if command == "pwd" {
			// The os.Getwd function is used to get the current working directory.
			// If the function returns an error, the error message is printed to the console.
			// Otherwise, the current working directory is printed.
			if wd, err := os.Getwd(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(wd)
			}
			continue
		}

		if command == "cd" {
			// The os.Chdir function is used to change the current working directory.
			// If the function returns an error, the error message is printed to the console.
			// Otherwise, the current working directory is printed.
			if len(inputParts) < 2 {
				fmt.Println("cd: missing argument")
				continue
			}
			commandArg := inputParts[1]
			if err := os.Chdir(commandArg); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", commandArg)
			}
			continue
		}

		// Execute external command
		cmd := exec.Command(command, inputParts[1:]...)

		// Set up command output and error handling
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Run the command
		err = cmd.Run()

		// DEFAULT
		// prints a message indicating that the command is not found.
		// The input[:len(input)-1] part removes the newline character from the end of the input string, ensuring the command name is printed correctly without an extra line break.
		// This change allows the shell to handle invalid commands by displaying a message in the format <command_name>: command not found
		if err != nil {
			fmt.Printf("%s: command not found\n", input[:])
		}
	}
}
