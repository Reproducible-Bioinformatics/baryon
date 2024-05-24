package parser

import (
	"baryon/tool"
	"fmt"
	"regexp"
	"strings"
)

// roxygen implements the functions to parse R function documentation
// and obtain a Galaxy Tool.
type roxygen struct{}

// NewRoxygen returns a New roxygen.
func NewRoxygen() *roxygen {
	return &roxygen{}
}

func (*roxygen) Parse(in []byte) (*tool.Tool, error) {
	var outtool tool.Tool
	comment := obtainComment(in)
	if len(comment) == 0 {
		return nil, fmt.Errorf("Cannot parse roxygen comment.")
	}
	return &outtool, nil
}

var commentEntryRegex = regexp.MustCompile(`@[^@]+`)

func getCommentEntries(input string) []string {
	return commentEntryRegex.FindAllString(input, -1)
}

// Obtains the roxygen comment form the input "in".
func obtainComment(in []byte) string {
	var commentLines string
	for _, line := range strings.Split(string(in), "\n") {
		if ok, submatches := isRoxygenLine(line); ok {
			for _, s := range submatches[1:] {
				commentLines += s + "\n"
			}
		}
	}
	return commentLines
}

// roxygenLineRegex, matches a roxygenline.
var roxygenLineRegex = regexp.MustCompile("^#' ?(.*)")

// Returns true if a line is a roxygen line, false otherwise.
func isRoxygenLine(line string) (bool, []string) {
	submatches := roxygenLineRegex.FindStringSubmatch(line)
	return len(submatches) > 0, submatches
}
