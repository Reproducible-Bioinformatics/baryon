package parser

import (
	"baryon/tool"
	"regexp"
)

type Parser interface {
	// Parse parses a []byte.
	Parse([]byte) (*tool.Tool, error)
}

// roxygen implements the functions to parse R function documentation
// and obtain a Galaxy Tool.
type roxygen struct{}

// NewRoxygen returns a New roxygen.
func NewRoxygen() *roxygen {
	return &roxygen{}
}

func (*roxygen) Parse(in []byte) (*tool.Tool, error) {
	return &tool.Tool{}, nil
}

// roxygenLineRegex, matches a roxygenline.
var roxygenLineRegex = regexp.MustCompile("^#' ?.*")

// Returns true if a line is a roxygen line, false otherwise.
func isRoxygenLine(line string) bool {
	return roxygenLineRegex.MatchString(line)
}
