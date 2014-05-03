package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/jweslley/procker"
)

func main() {
	procfile := flag.String("f", "Procfile", "Procfile declaring commands to run")

	file, err := os.Open(*procfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	processes, err := procker.ParseProcfile(file)
	if err != nil {
		log.Fatalf("procker: %v", err)
	}
	log.Println(processes)

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
		process.Start(wd, []string{}, os.Stdout, os.Stderr)
	}

	for _, process := range processes {
		process.Wait()
	}
}
