package main

import "fmt"

var cmdHelp = &command{
	desc: "Show this help",
	help: `Usage: procker help [command]

Help shows usage for a command.`,
	call: help}

func help(args []string) {
	fmt.Printf("Help %s version %s\n", programName, programVersion)
}
