package builtins

import (
	. "fmt"
	"os"
	"errors"
)

func Cd(args []string) (err error) {

	/* "~" shortcut to /home/user directory */
	if args[1] == "~" || args[1] == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			args[1] = homeDir
		}
	}

	os.Chdir(args[1])

	return nil
}

func Help(args []string) (err error) {

	Print("goshell: simple shell written in go\n" +
		  "Type program names and arguments, and hit enter.\n" +
		  "Use the man command for informations about programs.\n")

	return nil
}

func Exit(args []string) (err error) {
	return errors.New("exit")
}
