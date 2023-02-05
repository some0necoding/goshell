package main

import (
	"bufio"
	"errors"
	. "fmt"
	"os"
	"os/exec"
	"shell/config"
	"strings"
)

/* Custom type for builtin functions */
type builtInFunc func([]string) error

/* Builtin functions' names */
var builtInStr = []string{
	"cd",
	"help",
	"exit",
}

/* Builtin functions */
var builtInFuncs = []builtInFunc{
	cd,
	help,
	exit,
}

func main() {

	// config files

	loop()

	// shutdown / cleanup
}

func loop() {

	var status error
	prompt := config.Prompt()

	for status == nil {
	
		Print(prompt)
		
		line, err := readLine()

		if err == nil {
			args := strings.Split(line, " ")
			status = execute(args)
		} else {
			status = err
		}
	}
}

func readLine() (line string, err error) {

	stdin := bufio.NewReader(os.Stdin)
	
	if userInput, err := stdin.ReadString('\n'); err == nil {
		return strings.TrimRight(userInput, "\n"), nil
	}
	
	return "", err
}

func execute(args []string) (err error) {

	// checking for empty commands
	if args[0] == "" {
		return nil
	}

	// checking for built-in commands
	for i, str := range builtInStr {
		if args[0] == str {
			function := builtInFuncs[i]
			return function(args)
		}
	}

	// otherwise launch the non built-in command
	return launch(args)
}

// launches non built-in commands
func launch(args []string) (err error) {

	// starts the command and waits for it to finish
	if process, err := Start(args); err == nil {
		process.Wait()
	}

	return nil
}

/* Start a new process */
func Start(args []string) (process *os.Process, err error) {

	/* Check if command exists */
	if args[0], err = exec.LookPath(args[0]); err == nil {

		/* Start the new process */
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
	BUILT-IN FUNCTIONS DECLARATION START
*/

// built-in cd command
func cd(args []string) (err error) {

	// "~" shortcut to /home/user directory
	if args[1] == "~" || args[1] == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			args[1] = homeDir
		}
	}

	os.Chdir(args[1])

	return nil
}

// built-in help command
func help(args []string) (err error) {

	Print("goshell: simple shell written in go\n" +
		  "Type program names and arguments, and hit enter\n" +
		  "\nThe following commands are built-in:\n")

	for _, str := range builtInStr {
		Println("\t", str)
	}

	Println("\nUse the man command for informations about other programs.")
	return nil
}

// built-in exit command
func exit(args []string) (err error) {
	return errors.New("exit")
}

/*
 BUILT-IN FUNCTIONS DECLARATION END
*/
