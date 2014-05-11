package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/jweslley/procker"
)

const defaultEnvfile = ".env"

var (
	cmdStart = &command{
		desc: "Start application's processes",
		help: `Usage: procker start [options]

Start the processes specified by a Procfile

Available options:`,
		exec: start,
		flag: startFlags}

	// flags
	startFlags    = flag.NewFlagSet("start", flag.ExitOnError)
	startProcfile = startFlags.String("f", "Procfile",
		"Procfile declaring commands to run")
	startEnvfile = startFlags.String("e", defaultEnvfile,
		"File containing environment variables to be used")
	startBasePort = startFlags.Int("p", 5000,
		"Base port to be used by processes. Should be a multiple of 1000")
)

func start(args []string) {
	procSpecs := parseProfile(*startProcfile)
	env := parseEnv(*startEnvfile)
	dir := path.Dir(*startProcfile)
	padding := longestName(procSpecs)
	log.SetOutput(procker.NewPrefixedWriter(os.Stdout, prefix(programName, padding)))
	process := buildProcess(procSpecs, dir, env, *startBasePort, padding)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("%v received, stopping processes and exiting.", sig)
			process.Kill()
			os.Exit(1)
		}
	}()

	err := process.Start()
	failIf(err)

	err = process.Wait()
	failIf(err)
}

func buildProcess(
	specs map[string]string,
	dir string,
	env []string,
	basePort int,
	padding int) procker.Process {

	i := 0
	p := []procker.Process{}
	for name, command := range specs {
		port := basePort + (i * 100)
		process := procker.NewProcess(name,
			command,
			dir,
			append(env, fmt.Sprintf("PORT=%d", port)),
			procker.NewPrefixedWriter(os.Stdout, prefix(name, padding)),
			procker.NewPrefixedWriter(os.Stderr, prefix(name, padding)))
		p = append(p, process)
		i++

		log.Printf("starting %s on port %d", name, port)
	}
	return procker.NewProcessSet(p...)
}

func parseProfile(filepath string) map[string]string {
	file, err := os.Open(filepath)
	defer file.Close()
	failIf(err)

	processes, err := procker.ParseProcfile(file)
	failIf(err)
	return processes
}

func parseEnv(filepath string) []string {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		if filepath == defaultEnvfile {
			return []string{}
		} else {
			failIf(err)
		}
	}

	env, err := procker.ParseEnv(file)
	failIf(err)
	return env
}

func longestName(processes map[string]string) int {
	max := len(programName)
	for name := range processes {
		if len(name) > max {
			max = len(name)
		}
	}
	return max
}

func prefix(prefix string, padding int) string {
	return fmt.Sprintf(fmt.Sprintf("%%%ds | ", -padding), prefix)
}
