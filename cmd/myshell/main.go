package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		// reads a line of input from the user and stores it in the variable input.
		// The underscore _ is used to ignore any error returned by ReadString.
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
				if _, err := os.Stat(path.Join("/bin", inputParts[1])); err == nil {
					fmt.Printf("%s is /bin/%s\n", inputParts[1], inputParts[1])
				} else if _, err := os.Stat(path.Join("/usr/bin", inputParts[1])); err == nil {
					fmt.Printf("%s is /usr/bin/%s\n", inputParts[1], inputParts[1])
				} else {
					fmt.Printf("%s: not found\n", inputParts[1])
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
