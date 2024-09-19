package main

import (
	"baryon/marshaler"
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
	parser := parser.NewRoxygen()
	file, err := getFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if stat.Size() == 0 {
		log.Fatal("No file provided.")
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
	output, err := func() ([]byte, error) {
		var mode string
		if len(argsWithoutProg) > 1 {
			mode = argsWithoutProg[1]
		}
		switch mode {
		case "bash":
			out, err := marshaler.BashMarshaler{}.Marshal(tool)
			return out, err
		default:
			out, err := xml.MarshalIndent(tool, "", "\t")
			return out, err
		}
	}()
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
