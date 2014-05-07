package main

import "fmt"

var cmdVersion = &command{
	desc: "Display current version",
	help: `Usage: procker version

Display current version`,
	call: version}

func version(args []string) {
	fmt.Printf("%s version %s\n", programName, programVersion)
}
