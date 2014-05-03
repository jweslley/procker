package main

import (
	"flag"
	"log"
	"os"
	"path"

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

	wd := path.Dir(*procfile)
	for name, process := range processes {
		log.Printf("starting %s - %s", name, process.Command)
		process.Start(wd, []string{}, os.Stdout, os.Stderr)
		process.Wait()
	}
}
