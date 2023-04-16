package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/surface-security/scanner-go-entrypoint/scanner"
)

type DomainListInput struct {
	Name    string
	Domains []string
}

func main() {
	s := scanner.Scanner{Name: "subfinder"}
	options := s.BuildOptions()
	scanner.ParseOptions(options)

	err := os.MkdirAll(options.Output, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	scanner.ReadInputJSONLines(options, func(input DomainListInput) bool {
		if len(input.Domains) == 0 {
			log.Printf("No domains for %s", input.Name)
			return true
		}
		// pass temporary file to subfinder instead of final path, as only finished files should be placed there
		file, err := os.CreateTemp("", "subfinder")
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer os.Remove(file.Name())

		err = s.Exec(
			"-json",
			"-o", file.Name(),
			// no point checking for updates
			"-duc",
			"-d", strings.Join(input.Domains, ","),
		)
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
		return true
	})
}
