package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/jweslley/procker"
)

const defaultEnvfile = ".env"

var (
	cmdStart = &command{
		desc: "Start application's processes",
		help: `Usage: procker start [options] [process name]...

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
		"Base port to be used by processes")
	startStopTimeout = startFlags.Int("t", 5,
		"Time (in seconds) for graceful stop of processes")
)

func start(args []string) {
	processes := parseProfile(*startProcfile)
	env := parseEnv(*startEnvfile)
	dir := path.Dir(*startProcfile)
	padding := longestName(processes)
	log.SetFlags(0)
	log.SetOutput(procker.NewPrefixedWriter(os.Stdout, prefix(programName, padding)))
	process := buildProcess(args, processes, dir, env, *startBasePort, padding)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		stopping := false
		for sig := range c {
			if stopping {
				log.Printf("%v signal received, killing processes and exiting.", sig)
				process.Signal(syscall.SIGKILL)
			} else {
				log.Printf("%v signal received, stopping processes and exiting.", sig)
				stopping = true
				go func() {
					process.Stop(time.Duration(*startStopTimeout) * time.Second)
				}()
			}
		}
	}()

	err := process.Start()
	failIf(err)

	process.Wait()
}

func buildProcess(
	processNames []string,
	processes map[string]string,
	dir string,
	env []string,
	port, padding int) procker.Process {

	p := []procker.Process{}
	for name, command := range processes {
		if !mustStart(processNames, name) {
			continue
		}

		process := &procker.SysProcess{
			Command:     command,
			Dir:         dir,
			Env:         append(env, fmt.Sprintf("PORT=%d", port)),
			Stdout:      procker.NewPrefixedWriter(os.Stdout, prefix(name, padding)),
			Stderr:      procker.NewPrefixedWriter(os.Stderr, prefix(name, padding)),
			SysProcAttr: sysProcAttrs(),
		}

		log.Printf("starting %s on port %d", name, port)
		p = append(p, process)
		port++
	}

	if len(p) == 0 {
		fail("no process to run\n")
	}

	return procker.NewProcessGroup(p...)
}

func mustStart(processNames []string, name string) bool {
	if len(processNames) == 0 {
		return true
	}

	for _, process := range processNames {
		if process == name {
			return true
		}
	}
	return false
}

func parseProfile(filepath string) map[string]string {
	file, err := os.Open(filepath)
	failIf(err)
	defer file.Close()

	processes, err := procker.ParseProcfile(file)
	failIf(err)
	return processes
}

func parseEnv(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		if filepath == defaultEnvfile {
			return os.Environ()
		} else {
			failIf(err)
		}
	}
	defer file.Close()

	env, err := procker.ParseEnv(file)
	failIf(err)
	return append(os.Environ(), env...)
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
