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

const programName = "procker"

func main() {
	procfile := flag.String("f", "Procfile", "Procfile declaring commands to run")
	envfile := flag.String("e", ".env", "File containing environment variables to be used")
	flag.Parse()

	processes := parseProfile(*procfile)
	env := parseEnv(*envfile)

	padding := longestName(processes)
	log.SetOutput(procker.NewPrefixedWriter(os.Stdout, prefix(programName, padding)))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("%v received, stopping processes and exiting.", sig)
			for name, process := range processes {
				log.Printf("killing %s", name)
				process.Kill()
			}
			os.Exit(1)
		}
	}()

	wd := path.Dir(*procfile)
	for name, process := range processes {
		log.Printf("starting %s - %s", name, process.Command)
		process.Start(wd,
			env,
			procker.NewPrefixedWriter(os.Stdout, prefix(name, padding)),
			procker.NewPrefixedWriter(os.Stderr, prefix(name, padding)))
	}

	for _, process := range processes {
		process.Wait()
	}
}

func parseProfile(filepath string) map[string]*procker.Process {
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

func longestName(processes map[string]*procker.Process) int {
	max := len(programName)
	for name, _ := range processes {
		if len(name) > max {
			max = len(name)
		}
	}
	return max
}

func prefix(prefix string, padding int) string {
	return fmt.Sprintf(fmt.Sprintf("%%%ds | ", -padding), prefix)
}
