package main

import (
	"fmt"
	"os"
)

var cmdHelp = &command{
	desc: "Show this help",
	help: `Usage: procker help [command]

Help shows usage for a command.`,
	call: help}

func help(args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(0)
	}

	command := findCommand(args[0])
	fmt.Println(command.help)
}
