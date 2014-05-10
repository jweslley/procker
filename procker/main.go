package main

import (
	"flag"
	"fmt"
	"os"
)

const programName = "procker"
const programVersion = "0.0.1"

type command struct {
	desc string
	help string
	exec func([]string)
	flag *flag.FlagSet
}

var commands map[string]*command

func init() {
	commands = map[string]*command{
		"start":   cmdStart,
		"version": cmdVersion,
		"help":    cmdHelp}
}

func main() {
	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	command := findCommand(os.Args[1])
	command.flag.Parse(os.Args[2:])
	command.exec(command.flag.Args())
}

func findCommand(name string) *command {
	c, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "procker: '%s' is not a procker command. See 'procker help'.\n", name)
		os.Exit(1)
	}
	if c.flag == nil {
		c.flag = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return c
}

func usage() {
	fmt.Println(`Usage: procker <command> [<args>]

Available commands:`)
	for name, command := range commands {
		fmt.Printf("%10s  %s\n", name, command.desc)
	}
	fmt.Println("\nRun 'procker help [command]' for details.")
}
