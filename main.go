package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	input      string
	output     string
	binPath    string
	extraHelp  bool
	extraFlags []string
}

// BuildOptions parses the command line flags provided by a user
func BuildOptions() *Options {
	options := &Options{}
	flag.StringVarP(&options.output, "output", "o", "/output", "Scanner results directory")
	flag.StringVarP(&options.binPath, "bin", "b", "subfinder", "Path to scanner binary")
	flag.BoolVarP(&options.extraHelp, "scanner-help", "H", false, "Show help for the scanner extra flags")
	return options
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions(options *Options) {
	flag.Parse()

	if flag.CommandLine.NArg() > 0 {
		args := flag.CommandLine.Args()
		options.extraFlags = args[:len(args)-1]
		options.input = args[len(args)-1]
	}
}

type SurfaceBugBountyInput struct {
	Name    string
	Domains []string
}

func main() {
	// Parse the command line flags and read config files
	options := BuildOptions()
	ParseOptions(options)

	if options.extraHelp {
		cmd := exec.Command(options.binPath, "-h")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to run scanner: %v", err)
		}
		exe := os.Args[0]
		fmt.Println(`
## Note ##
In order to pass any of these flags to the scanner, append them to the end of the command line, after "--".

Normal: ` + exe + ` ... /path/to/input.txt
Extra flags: ` + exe + ` ... -- -extra -flags /path/to/input.txt`)
		// same exit code as normal help
		os.Exit(2)
	}
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
		var input SurfaceBugBountyInput

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

		flags := append(
			[]string{
				"-json",
				"-o", file.Name(),
				// no point checking for updates
				"-duc",
				"-d", strings.Join(input.Domains, ","),
			},
			options.extraFlags...,
		)
		cmd := exec.Command(options.binPath, flags...)
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
