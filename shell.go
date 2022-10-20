package main

import (
	"bufio"
	. "fmt"
	"os"
	"strings"
	"syscall"
)

type builtInFunc func([]string) int

var builtInStr = []string{
	"cd",
	"help",
	"exit",
}

var builtInFuncs = []builtInFunc{
	cd,
	help,
	exit,
}

func main() {

	// load config files

	loop()

	// shutdown / cleanup
}

func loop() {

	status := 0

	for status == 0 {
		Print("> ")
		line, err := readLine()

		if err == 0 {
			args := splitLine(line)
			status = execute(args)
		}
	}
}

func readLine() (string, int) {

	var userInput string

	stdin := bufio.NewReader(os.Stdin)
	userInput, err := stdin.ReadString('\n')

	if err != nil {
		return "", -1
	}

	lastIndex := len(userInput) - 1

	return userInput[:lastIndex], 0
}

func splitLine(line string) []string {
	return strings.Split(line, " ")
}

func execute(args []string) int {

	if args[0] == "" {
		return -1
	}

	for i := 0; i < numBuiltins(); i++ {
		if args[0] == builtInStr[i] {
			function := builtInFuncs[i]
			return function(args)
		}
	}

	return launch(args)
}

func numBuiltins() int {
	return len(builtInStr)
}

func launch(args []string) int {

	pid, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)

	if pid == 0 {

		err := syscall.Exec(args[0], args, os.Environ())

		if err != nil {
			return -1
		}

	} else if pid < 0 {
		return -1
	} else {

		var status syscall.WaitStatus

		for !status.Exited() && !status.Signaled() {

			process, err := os.FindProcess(int(pid))

			if err == nil {

				processState, err := process.Wait()

				if err == nil {
					status = processState.Sys().(syscall.WaitStatus)
				} else {
					return -1
				}

			} else {
				return -1
			}
		}
	}

	return 0
}

/*
	BUILT-IN FUNCTION DECLARATION START
*/

func cd(args []string) int {
	if args[1] == "" {
		Println("ghs: expected argument to \"cd\"")
		return 1
	} else {
		os.Chdir(args[1])
	}

	return 0
}

func help(args []string) int {

	Print("GSH: a simple shell written in go\n" +
		"Type program names and arguments, and hit enter\n" +
		"\nThe following are built-in:\n")

	for i := 0; i < numBuiltins(); i++ {
		Println("\t", builtInStr[i])
	}

	Println("\nUse the man command for informations about other programs.")
	return 0
}

func exit(args []string) int {
	return 1
}

/*
 BUILT-IN FUNCTION DECLARATION END
*/
