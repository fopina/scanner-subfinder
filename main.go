package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	input   string
	output  string
	binPath string
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}
	flag.StringVarP(&options.output, "output", "o", "/output", "Scanner results directory")
	flag.StringVarP(&options.binPath, "bin", "b", "subfinder", "Path to scanner binary")

	flag.Parse()

	if flag.CommandLine.NArg() > 0 {
		options.input = flag.CommandLine.Arg(0)
	}
	return options

}

type SurfaceInput struct {
	Name    string
	Domains []string
}

func main() {
	// Parse the command line flags and read config files
	options := ParseOptions()
	err := os.MkdirAll(options.output, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}
	jsonFile, err := os.Open(options.input)
	if err != nil {
		log.Fatalf("%v", err)
	}
	dec := json.NewDecoder(jsonFile)
	for {
		var input SurfaceInput

		err := dec.Decode(&input)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}

		// pass temporary file to subfinder instead of final path, as only finished files should be placed there
		file, err := os.CreateTemp("", "subfinder")
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer os.Remove(file.Name())

		cmd := exec.Command(options.binPath, "-json", "-o", file.Name(), "-d", strings.Join(input.Domains, ","))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			log.Fatalf("Failed to run scanner: %v", err)
		}

		realOutputFile := path.Join(options.output, input.Name)
		outputFile, err := os.Create(realOutputFile)
		if err != nil {
			log.Fatalf("Couldn't open dest file: %v", err)
		}
		defer outputFile.Close()
		_, err = io.Copy(outputFile, file)
		if err != nil {
			log.Fatalf("Writing to output file failed: %v", err)
		}
	}
}
