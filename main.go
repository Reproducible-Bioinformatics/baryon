package main

import (
	"log"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var filePath string
	if len(argsWithoutProg) > 0 {
		filePath = argsWithoutProg[0]
	}
	file, err := getFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

// getFile retrieves a *os.File if a path is provided and is not empty.
// Otherwise, it obtains os.Stdin.
func getFile(path string) (*os.File, error) {
	if path == "" {
		return os.Stdin, nil
	}
	return os.Open(path)
}
