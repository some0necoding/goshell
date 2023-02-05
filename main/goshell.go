package main

import (
	"bufio"
	. "fmt"
	"os"
	"os/exec"
	"shell/config"
	"shell/builtins"
	"strings"
)

/* Custom type for builtin functions */
type builtInFunc func([]string) error

var builtInFuncs = map[string]builtInFunc {
	"cd": builtins.Cd,
	"help": builtins.Help,
	"exit": builtins.Exit, 
}

func main() {

	/* Config files */

	loop()

	/* Shutdown / Cleanup */
}

func loop() {

	var status error
	prompt := config.Prompt()

	for status == nil {

		Print(prompt)
		line, err := readLine()
		status = err

		if err == nil {
			args := strings.Split(line, " ")
			status = execute(args)
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

	/* Check for empty commands */
	if args[0] == "" {
		return nil
	}

	/* Check for built-in commands */
	if function, ok := builtInFuncs[args[0]]; ok {
		return function(args)
	}

	/* Otherwise launch non built-in command */
	return launch(args)
}

/* This function launches non built-in commands */
func launch(args []string) (err error) {

	/* Start the process and wait until finish */
	if process, err := Start(args); err == nil {
		process.Wait()
	}

	return err
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
