package main

import (
	"bufio"
	"errors"
	. "fmt"
	"os"
	"os/exec"
	"shell/builtins"
	"shell/config"
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

func loop() (status error) {

	prompt := config.Prompt()

	for status == nil {

		var line string

		Print(prompt)
		line, err := readLine()

		if err != nil {
			continue
		}

		args, err := splitLine(line)

		if err != nil {
			continue
		}
			
		status = execute(args)
	}

	return
}

func readLine() (line string, err error) {

	stdin := bufio.NewReader(os.Stdin)
	
	if userInput, err := stdin.ReadString('\n'); err == nil {
		return strings.TrimRight(userInput, "\n"), nil
	}
	
	return "", err
}

func splitLine(line string) (args []string, err error) {

	if strings.Count(line, "\"") % 2 == 1 {
		Println("gsh: missing \"")
		return args, errors.New("Bad formatting")
	} 

	tokens := strings.Split(line, "\"")

	for i, token := range tokens {
		
		token = strings.TrimSpace(token)

		if i % 2 == 1 {
			args = append(args, token)
		} else {
			args = append(args, strings.Split(token, " ")...)
		}
	}

	return args, nil
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
		_, err = process.Wait()
	}

	return err
}

/* Start a new process */
func Start(args []string) (process *os.Process, err error) {

	/* Check if command exists */
	if args[0], err = exec.LookPath(args[0]); err == nil {

		/* Start the new process */
		process, err := os.StartProcess(args[0], args, &os.ProcAttr{
			/* Set process files (same as parent) */
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
