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
	commentEntries := getCommentEntries(comment)
	for _, commentEntry := range commentEntries {
		split := strings.Split(commentEntry, " ")
		if len(split) > 0 {
			keyword := strings.TrimLeft(split[0], "@")
			if matcher, ok := act[keyword]; ok {
				matcher(strings.Join(split[1:], " "), &outtool)
			}
		}
	}
	return &outtool, nil
}

type Actor func(string, *tool.Tool)

// act serves as the entrypoint to parse a roxygen comment entry.
// it provides a set of functions, parsing each field.
// Implementation is dependent on the field.
var act map[string]Actor = map[string]Actor{
	"param": func(s string, t *tool.Tool) {
		if t.Inputs == nil {
			t.Inputs = &tool.Inputs{}
		}
		// TODO: parse all comment
		splitString := strings.Split(s, " ")
		if len(splitString) == 0 {
			return
		}
		// Processing of the name variable according to Galaxy's specs.
		name := splitString[0]
		name = strings.Replace(name, ".", "__", -1) // Replaces "." with "__".

		help := strings.Join(splitString[1:], " ")

		// Start processing baryon instructions.
		baryonInstruction := baryonNamespaceRegex.FindString(help)

		// Processing of the help string according to Galaxy's specs.
		help = strings.Replace(help, baryonInstruction, "", -1)
		help = strings.TrimSpace(help)

		t.Inputs.Param = append(t.Inputs.Param, tool.Param{
			// First element is the name of the param.
			Name: name,
			// Help is the other part of the string.
			Help: help,
		})
	},
	"description": func(content string, t *tool.Tool) {
		t.Description = content
	},
	"author": func(content string, t *tool.Tool) {
		for _, name := range strings.Split(content, ",") {
			if t.Creator == nil {
				t.Creator = &tool.Creator{}
			}
			t.Creator.Person = append(
				t.Creator.Person, tool.Person{
					Name: strings.TrimSpace(name),
				})
		}
	},
}

var baryonNamespaceRegex = regexp.MustCompile(`\$B{([^}]*)}`)

// Regex to obtain a comment entry.
var commentEntryRegex = regexp.MustCompile(`@[^@]+`)

// Get all entries from a comment.
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
