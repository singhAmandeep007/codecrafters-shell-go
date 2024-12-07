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

func parseInput(s string) []string {
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

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

/**
 * main is the entry point of the shell program. It runs an infinite loop that continuously prompts the user for input,
 * processes the input, and executes the corresponding commands. The shell supports built-in commands such as "echo",
 * "type", "pwd", and "cd", as well as external commands found in the system's PATH.
 *
 * The main function performs the following steps:
 * 1. Prompts the user for input by printing "$ " to the standard output.
 * 2. Reads a line of input from the user and trims any leading or trailing whitespace.
 * 3. Splits the input string into a slice of strings, where the first element is the command and the rest are arguments.
 * 4. Checks if the command is a built-in command and executes the corresponding logic:
 *    - "echo": Prints the arguments to the console.
 *    - "type": Displays information about the specified command.
 *    - "pwd": Prints the current working directory.
 *    - "cd": Changes the current working directory to the specified path.
 * 5. If the command is not a built-in command, it attempts to execute it as an external command using the exec.Command function.
 * 6. If the command is not found, it prints an error message indicating that the command is not found.
 *
 * The shell continues to run until the user enters the "exit 0" command, which terminates the program.
 */
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
		inputParts := parseInput(input)

		command := strings.ToLower(inputParts[0])

		// If the user enters the exit command, the shell will exit.
		if input == "exit 0" {
			os.Exit(0)
		}

		if command == "echo" {
			commandArg := inputParts[1:]

			if len(commandArg) == 0 {
				fmt.Println("echo: missing argument")
				continue
			}
			// The first word is the command name, and the rest of the words are the arguments.
			// The arguments are joined together with a space character and printed to the console.
			fmt.Printf("%s\n", strings.Join(commandArg, " "))
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
			/**
			 * The added function cd takes an array of strings as its argument, which represents the command-line arguments passed to the cd command. Inside the function:
			 * The first argument, which should be the directory path, is accessed with inputParts[1].
			 * The os.Chdir function is called with this path to attempt to change the current working directory.
			 * If os.Chdir returns an error (which happens if the directory does not exist or cannot be accessed), an error message is printed to standard output using fmt.Printf. The message follows the format "<directory>: No such file or directory\n", where <directory> is replaced with the actual directory path that was attempted.
			 * Absolute paths, like /usr/local/bin
			 * Relative paths, like ./, ../, ./dir.
			 * Paths with special characters, like ~ which stands for the user's home directory. The home directory is specified by the HOME environment variable.
			 */

			// Handle special characters in the directory path
			if commandArg == "~" {
				// The os.Getenv function is used to retrieve the value of the HOME environment variable.
				// If the commandArg is equal to ~, the HOME environment variable is used as the directory path.
				commandArg = os.Getenv("HOME")
				// If the commandArg starts with ~/ (indicating a relative path from the home directory), the path is constructed by joining the home directory path with the rest of the commandArg.
				// EXAMPLE: If the HOME environment variable is /home/user and the commandArg is ~/Documents, the resulting path will be /home/user/Documents.
			} else if strings.HasPrefix(commandArg, "~/") {
				commandArg = filepath.Join(os.Getenv("HOME"), commandArg[2:])
			}

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
