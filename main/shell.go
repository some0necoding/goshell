package main

// #include "../cfuncs/launch.h"
import "C"

import (
	. "fmt"
	"strings"
)

type builtInFunc func([]string)

var builtInStr = []string {
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
		line := readLine()
		args := splitLine(line)
		status := execute(args)
	}
}

func readLine() string {

	var userInput string
	
	Scan(&userInput)

	return userInput
}

func splitLine(line string) []string {
	return strings.Split(line, " ") 
}

func execute(args []string) int {

	var c_args **C.char

	if args[0] == nil {
		return 1
	}

	for i := 0; i < numBuiltins(); i++ {
		if args[0] == builtInStr[i] {
			return builtInFuncs[i]
		}
	}

	c_args = toCStrArray(args)

	return C.launch(c_args)
}

func numBuiltins() int {
	return len(builtInStr)
}

func toCStrArray(strArr []string) **C.char {

	var cStrArr **C.char
	var cStr *C.char

	for i, str := range strArr {
		cStr := C.Cstring(str)
		cStrArr[i] = cStr
	}

	return cStrArr
}

/*
	BUILT-IN FUNCTION DECLARATION START
*/

func cd(args []string) int {
	if args[1] == nil {
		Println("ghs: expected argument to \"cd\"")
		return 1
	} else {
		// call to c chdir(args[1]) function
	}

	return 0
}

func help(args []string) int {

	Print("GSH: a simple shell written in go\n" +
		  "Type program names and arguments, and hit enter\n" +
		  "The following are built-in:\n")

	for i := 0; i < numBuiltins(); i++ {
		Println(builtinStr[i])
	}

	Println("Use the man command for informations about other programs.")
	return 0
}

func exit(args []string) int {
	return 0
}

/*
	 BUILT-IN FUNCTION DECLARATION END
*/