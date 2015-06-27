package main

import (
	"flag"
	"os"
	"strings"

	"github.com/jweslley/procker"
)

var (
	cmdRun = &command{
		desc: "Run a command using your application's environment",
		help: `Usage: procker run COMMAND

Run a command using your application's environment

Available options:`,
		exec: run,
		flag: runFlags}

	// flags
	runFlags   = flag.NewFlagSet("run", flag.ExitOnError)
	runEnvfile = runFlags.String("e", defaultEnvfile,
		"File containing environment variables to be used")
)

func run(args []string) {
	if len(args) == 0 {
		fail("you must specify a command. See 'procker help run'.\n")
	}

	env := parseEnv(*runEnvfile)
	command := strings.Join(args, " ")
	process := &procker.SysProcess{
		Command: command,
		Env:     env,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	err := process.Start()
	failIf(err)

	process.Wait()
}
