package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	// Attempts to increase the OS file descriptors - Fail silently
	_ "github.com/projectdiscovery/fdmax/autofdmax"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Options struct {
	runner.Options
	surfaceInput string
}

func createGroup(flagSet *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

// ConfigureOutput configures the output on the screen
func ConfigureOutput(options *Options) {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
	if options.NoColor {
		gologger.DefaultLogger.SetFormatter(formatter.NewCLI(true))
	}
	if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	flagSet, options := runner.BuildOptions()

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	newOptions := &Options{Options: *options}
	//flag.StringVar(&finalOutput, "o", "/output/output.txt", "Output results to file (Subjack will write JSON if file ends with '.json').")
	if flagSet.CommandLine.NArg() > 0 {
		newOptions.surfaceInput = flagSet.CommandLine.Arg(0)
	}

	// Default output is stdout
	options.Output = os.Stdout

	// Check if stdin pipe was given
	options.Stdin = false

	// Read the inputs and configure the logging
	ConfigureOutput(newOptions)

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", runner.Version)
		os.Exit(0)
	}

	options.AllSources = passive.DefaultAllSources
	options.Recursive = passive.DefaultRecursiveSources
	options.Recursive = resolve.DefaultResolvers
	options.Sources = passive.DefaultSources
	//options.Providers = &runner.Providers{}
	return newOptions
}

func runIt(options *Options) {
	newRunner, err := runner.NewRunner(&options.Options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	err = newRunner.RunEnumeration(context.Background())
	if err != nil {
		gologger.Fatal().Msgf("Could not run enumeration: %s\n", err)
	}
}

type SurfaceInput struct {
	Name    string
	Domains []string
}

func main() {
	// Parse the command line flags and read config files
	options := ParseOptions()
	if options.surfaceInput == "" {
		runIt(options)
	} else {
		// re-define here instead of being default so it doesn't break standard subfinder usage
		surfaceOutput := options.OutputDirectory
		if surfaceOutput == "" {
			surfaceOutput = "/output/"
		}
		err := os.MkdirAll(surfaceOutput, 0755)
		if err != nil {
			gologger.Fatal().Msgf("%v", err)
		}
		jsonFile, err := os.Open(options.surfaceInput)
		if err != nil {
			gologger.Fatal().Msgf("%v", err)
		}
		dec := json.NewDecoder(jsonFile)
		for {
			var input SurfaceInput

			err := dec.Decode(&input)
			if err == io.EOF {
				break
			}
			if err != nil {
				gologger.Fatal().Msgf("%v", err)
			}

			// pass temporary file to subfinder instead of final path, as only finished files should be placed there
			file, err := ioutil.TempFile("", "subfinder")
			if err != nil {
				gologger.Fatal().Msgf("%v", err)
			}
			defer os.Remove(file.Name())

			options.OutputFile = file.Name()
			options.JSON = true
			options.Domain = input.Domains

			runIt(options)

			realOutputFile := path.Join(surfaceOutput, input.Name)
			outputFile, err := os.Create(realOutputFile)
			if err != nil {
				gologger.Fatal().Msgf("Couldn't open dest file: %v", err)
			}
			defer outputFile.Close()
			_, err = io.Copy(outputFile, file)
			if err != nil {
				gologger.Fatal().Msgf("Writing to output file failed: %v", err)
			}
		}
	}
}
