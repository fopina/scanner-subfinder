package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fopina/scanner-go-entrypoint/scanner"
)

type SurfaceBugBountyInput struct {
	Name    string
	Domains []string
}

func main() {
	s := scanner.Scanner{
		Name: "subfinder",
	}
	options := s.BuildOptions()
	scanner.ParseOptions(options)

	err := os.MkdirAll(options.Output, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}
	jsonFile, err := os.Open(options.Input)
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
			options.ExtraFlags...,
		)
		cmd := exec.Command(options.BinPath, flags...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			log.Fatalf("Failed to run scanner: %v", err)
		}

		realOutputFile := path.Join(options.Output, input.Name)
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
