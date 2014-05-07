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

	cmdName := args[0]
	command, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "procker: '%s' is not a procker command. See 'procker help'.\n", cmdName)
		os.Exit(1)
	}

	fmt.Println(command.help)
}
