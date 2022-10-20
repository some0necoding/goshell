package config

import (
	"os"
	"strings"
)

func Prompt() string {

	prompt := "["
	user := os.Getenv("USER")

	if hostname, err := os.Hostname(); err == nil {
		prompt += hostname + "@"
	}

	prompt += user + " "

	if fullPath, err := os.Getwd(); err == nil {

		var workingDir string

		if homeDir, err := os.UserHomeDir(); err == nil && fullPath == homeDir {
			workingDir = "~"
		} else {
			dirs := strings.Split(fullPath, "/")
			workingDir = dirs[(len(dirs) - 1)]
		}

		prompt += workingDir
	}

	prompt += "]$ "

	return prompt
}