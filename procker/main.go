package main

import (
	"fmt"
	"os"
)

const programName = "procker"
const programVersion = "0.0.1"

type command struct {
	desc string
	help string
	call func([]string)
}

var commands = make(map[string]*command)

func init() {
	commands["start"] = cmdStart
	commands["version"] = cmdVersion
	commands["help"] = cmdHelp
}

func main() {
	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	cmdName := os.Args[1]
	command, ok := commands[cmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "procker: '%s' is not a procker command. See 'procker help'.\n", cmdName)
		os.Exit(1)
	}

	command.call(os.Args[2:])
}

func usage() {
	fmt.Println(`Usage: procker <command> [<args>]

Available commands:`)
	for name, command := range commands {
		fmt.Printf("%10s  %s\n", name, command.desc)
	}
	fmt.Println("\nRun 'procker help [command]' for details.")
}
