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

var (
	cmdStart = &command{
		desc: "Start processes",
		help: `Usage: procker start [options]

Start the processes specified by a Procfile

Available options:`,
		exec: start,
		flag: startFlags}

	// flags
	startFlags = flag.NewFlagSet("start", flag.ExitOnError)
	procfile   = startFlags.String("f", "Procfile",
		"Procfile declaring commands to run")
	envfile = startFlags.String("e", ".env",
		"File containing environment variables to be used")
	basePort = startFlags.Int("p", 5000,
		"Base port to be used by processes. Should be a multiple of 1000")
)

func start(args []string) {
	procSpecs := parseProfile(*procfile)
	env := parseEnv(*envfile)
	dir := path.Dir(*procfile)
	padding := longestName(procSpecs)
	process := buildProcess(procSpecs, dir, env, *basePort, padding)

	log.SetOutput(procker.NewPrefixedWriter(os.Stdout, prefix(programName, padding)))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("%v received, stopping processes and exiting.", sig)
			process.Kill()
			os.Exit(1)
		}
	}()

	log.Printf("starting processes")
	err := process.Start()
	log.Printf("error on start %s", err)

	log.Printf("waiting processes")
	err = process.Wait()
	log.Printf("error on wait %s", err)
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
		port := fmt.Sprintf("PORT=%d", basePort+(i*100))
		process := procker.NewProcess(name, command, dir, append(env, port),
			procker.NewPrefixedWriter(os.Stdout, prefix(name, padding)),
			procker.NewPrefixedWriter(os.Stderr, prefix(name, padding)))
		p = append(p, process)
		i += 1
	}
	return procker.NewProcessSet(p...)
}

func parseProfile(filepath string) map[string]string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("procker: %v", err)
	}
	defer file.Close()

	processes, err := procker.ParseProcfile(file)
	if err != nil {
		log.Fatalf("procker: %v", err)
	}
	return processes
}

func parseEnv(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		return []string{}
	}

	env, err := procker.ParseEnv(file)
	if err != nil {
		log.Fatalf("procker: %v", err)
	}
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
