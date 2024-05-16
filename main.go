package main

import (
	"baryon/parser"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var filePath string
	if len(argsWithoutProg) > 0 {
		filePath = argsWithoutProg[0]
	}
	parser := parser.NewRlangDocParser()
	file, err := getFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileread, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	tool, err := parser.Parse(fileread)
	if err != nil {
		log.Fatal(err)
	}
	output, err := xml.MarshalIndent(tool, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output))
}

// getFile retrieves a *os.File if a path is provided and is not empty.
// Otherwise, it obtains os.Stdin.
func getFile(path string) (*os.File, error) {
	if path == "" {
		return os.Stdin, nil
	}
	return os.Open(path)
}
