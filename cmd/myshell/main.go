package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {

	for {
		fmt.Print("$ ")

		// Wait for user input
		// reads a line of input from the user and stores it in the variable input.
		// The underscore _ is used to ignore any error returned by ReadString.
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		commands := strings.Split(input, " ")

		switch commands[0] {
		case "exit":
			// The exit command is implemented by checking if the first word of the input is "exit" and the second word is "0".
			code, err := strconv.Atoi(commands[1])
			if err != nil {
				// shell will exit with a non-zero exit code.
				os.Exit(1)
			}
			// shell will successfully exit, zero exit code.
			os.Exit(code)
		case "echo":
			// The echo command is implemented by splitting the input string into words using the strings.Split function.
			// The first word is the command name, and the rest of the words are the arguments.
			// The arguments are joined together with a space character and printed to the console.
			fmt.Print(strings.Join(strings.Split(input, " ")[1:], " "))
		default:
			// prints a message indicating that the command is not found.
			// The input[:len(input)-1] part removes the newline character from the end of the input string, ensuring the command name is printed correctly without an extra line break.
			// This change allows the shell to handle invalid commands by displaying a message in the format <command_name>: command not found
			fmt.Printf("%s: command not found\n", input[:len(input)-1])
		}
	}
}
