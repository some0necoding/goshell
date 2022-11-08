package main

import (
	"bufio"
	. "fmt"
	"os"
	"os/exec"
	"strings"
	"shell/config"
)

// custom type for builtin functions
type builtInFunc func([]string) int

// builtin functions' names
var builtInStr = []string{
	"cd",
	"help",
	"exit",
}

// builtin functions
var builtInFuncs = []builtInFunc{
	cd,
	help,
	exit,
}

var prompt string

func main() {

	// config files

	loop()

	// shutdown / cleanup
}

func loop() {

	status := 0

	for status == 0 {
	
		// generating prompt
		prompt = config.Prompt()

		Print(prompt)
		
		line, err := readLine()

		if err == 0 {
			args := splitLine(line)
			status = execute(args)
		}
	}
}

func readLine() (string, int) {

	stdin := bufio.NewReader(os.Stdin)
	
	if userInput, err := stdin.ReadString('\n'); err == nil {
		lastIndex := len(userInput)	- 1
		
		// returning user input without trailing "\n"
		return userInput[:lastIndex], 0
	}
	
	return "", -1
}

func splitLine(line string) []string {
	return strings.Split(line, " ")
}

func execute(args []string) int {

	// checking for empty commands
	if args[0] == "" {
		return 0
	}

	// checking for built-in commands
	for i := 0; i < numBuiltins(); i++ {
		if args[0] == builtInStr[i] {
			function := builtInFuncs[i]
			return function(args)
		}
	}

	// otherwise launch the non built-in command
	return launch(args)
}

// returns number of built-in commands
func numBuiltins() int {
	return len(builtInStr)
}

// launches non built-in commands
func launch(args []string) int {

	// starts the command and waits for it to finish
	if process, err := Start(args); err == nil {
		process.Wait()
	}

	return 0
}

// starts a new process
func Start(args []string) (process *os.Process, err error) {

	// checking for command existance
	if args[0], err = exec.LookPath(args[0]); err == nil {

		// starting the new process
		process, err := os.StartProcess(args[0], args, &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})

		if err == nil {
			return process, nil
		}

	} else {
		Println("gsh: command not found")
	}

	return nil, err
}

/*
	BUILT-IN FUNCTION DECLARATION START
*/

// built-in cd command
func cd(args []string) int {

	// checking for bad args
	if args[1] == "" {
		Println("gsh: expected argument to \"cd\"")
		return 1
	} else {

		// "~" shortcut to /home/user directory
		if args[1] == "~" {
			if homeDir, err := os.UserHomeDir(); err == nil {
				args[1] = homeDir
			}
		}

		// making system call
		os.Chdir(args[1])
	}

	return 0
}

// built-in help command
func help(args []string) int {

	Print("goshell: simple shell written in go\n" +
		  "Type program names and arguments, and hit enter\n" +
		  "\nThe following commands are built-in:\n")

	for i := 0; i < numBuiltins(); i++ {
		Println("\t", builtInStr[i])
	}

	Println("\nUse the man command for informations about other programs.")
	return 0
}

// built-in exit command
func exit(args []string) int {
	return 1
}

/*
 BUILT-IN FUNCTION DECLARATION END
*/
