package parser

import (
	"baryon/tool"
	"fmt"
	"log"
	"os"
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
			matcher, ok := act[keyword]
			if !ok {
				continue
			}
			err := matcher(strings.Join(split[1:], " "), &outtool)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return &outtool, nil
}

type Actor func(string, *tool.Tool) error

// act serves as the entrypoint to parse a roxygen comment entry.
// it provides a set of functions, parsing each field.
// Implementation is dependent on the field.
var act map[string]Actor = map[string]Actor{
	"param": func(s string, t *tool.Tool) error {
		if t.Inputs == nil {
			t.Inputs = &tool.Inputs{}
		}
		splitString := strings.Split(s, " ")
		if len(splitString) == 0 {
			return nil
		}
		// Processing of the name variable according to Galaxy's specs.
		name := splitString[0]
		name = strings.Replace(name, ".", "__", -1) // Replaces "." with "__".
		help := strings.Join(splitString[1:], " ")

		// Start processing baryon instructions.
		baryonInstruction := baryonNamespaceRegex.FindStringSubmatch(help)

		// Processing of the help string according to Galaxy's specs.
		if len(baryonInstruction) > 0 {
			help = strings.Replace(help, baryonInstruction[0], "", -1)
		}
		help = strings.TrimSpace(help)

		newParam := tool.Param{
			Name:     name,
			Help:     help,
			Optional: true,
		}
		// Matched inside Baryon namespace.
		if len(baryonInstruction) > 1 {
			parseInstruction(&newParam, baryonInstruction[1])
		}
		err := newParam.Validate()
		if err != nil {
			return fmt.Errorf("Error in act.param: %v", err)
		}
		t.Inputs.Param = append(t.Inputs.Param, newParam)
		return nil
	},
	"description": func(content string, t *tool.Tool) error {
		baryonInstruction :=
			baryonNamespaceRegex.FindStringSubmatch(content)
		// Processing of the description string according to Galaxy's specs.
		if len(baryonInstruction) > 0 {
			content =
				strings.Replace(content, baryonInstruction[0], "", -1)
		}
		content = strings.TrimSpace(content)
		if len(baryonInstruction) < 1 {
			return nil
		}
		t.Description = content
		instructions := strings.Split(baryonInstruction[1], ";")
		if t.Requirements == nil {
			t.Requirements = &tool.Requirements{}
		}
		for _, instruction := range instructions {
			instruction = strings.TrimSpace(instruction)
			match := instructionRegex.FindStringSubmatch(instruction)
			if len(match) < 3 {
				continue
			}
			parser, ok := descriptionInstruction[match[1]]
			if !ok {
				fmt.Fprintf(os.Stderr,
					"Instruction \"%s\" not found. Continuing.\n", match[1])
				continue
			}
			parser(t, match[2])
		}
		return nil
	},
	"author": func(content string, t *tool.Tool) error {
		for _, name := range strings.Split(content, ",") {
			if t.Creator == nil {
				t.Creator = &tool.Creator{}
			}
			t.Creator.Person = append(
				t.Creator.Person, tool.Person{
					Name: strings.TrimSpace(name),
				})
		}
		return nil
	},
	"return": func(description string, t *tool.Tool) error {
		baryonInstruction :=
			baryonNamespaceRegex.FindStringSubmatch(description)
		// Processing of the description string according to Galaxy's specs.
		if len(baryonInstruction) > 0 {
			description =
				strings.Replace(description, baryonInstruction[0], "", -1)
		}
		description = strings.TrimSpace(description)
		if len(baryonInstruction) < 1 {
			return nil
		}
		instructions := strings.Split(baryonInstruction[1], ";")
		if t.Outputs == nil {
			t.Outputs = &tool.Outputs{}
		}
		for _, instruction := range instructions {
			instruction = strings.TrimSpace(instruction)
			match := instructionRegex.FindStringSubmatch(instruction)
			if len(match) < 3 {
				continue
			}
			parser, ok := returnInstructions[match[1]]
			if !ok {
				fmt.Fprintf(os.Stderr,
					"Instruction \"%s\" not found. Continuing.\n", match[1])
				continue
			}
			parser(t.Outputs, match[2])
		}
		return nil
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

// ParamFunction used to provide functions for Baryon Namespaces used inside
// roxygen2 params.
type ParamFunction func(*tool.Param, string)

// paramIstructions is a map of function used to when parsing a roxygen2 param.
var paramIstructions map[string]ParamFunction = map[string]ParamFunction{
	"!":        func(t *tool.Param, arg string) { t.Optional = false },
	"required": func(t *tool.Param, arg string) { t.Optional = false },
	"type":     func(t *tool.Param, arg string) { t.Type = arg },
	"value":    func(t *tool.Param, arg string) { t.Value = arg },
	"options": func(t *tool.Param, arg string) {
		for _, entry := range strings.Split(arg, ",") {
			trimmedSpace := strings.TrimSpace(entry)
			if trimmedSpace == "" {
				continue
			}
			t.Options = append(t.Options, tool.Option{
				Value:         trimmedSpace,
				CanonicalName: trimmedSpace, // TODO: Issue #4.
			})
		}
	},
}

// descriptionInstruction is a map of functions used when parsing roxygen2 return.
var descriptionInstruction map[string]DescriptionFunction = map[string]DescriptionFunction{
	"container": func(t *tool.Tool, args string) {
		argList := strings.Split(args, ",")
		if len(argList) < 1 {
			return
		}
		container := tool.Container{
			Type:  "docker", // This is the default for baryon.
			Value: strings.TrimSpace(argList[0]),
		}
		if len(argList) > 1 {
			container.Type = strings.TrimSpace(argList[1])
		}
		if err := container.Validate(); err != nil {
			log.Fatalf("Error in descriptionInstruction, container: %v", err)
		}
		if t.Requirements == nil {
			t.Requirements = &tool.Requirements{}
		}
		t.Requirements.Container = append(t.Requirements.Container, container)
	},
}

// ReturnFunction used to provide functions for Baryon Namespaces used inside
// roxygen2 return.
type DescriptionFunction func(*tool.Tool, string)

// ReturnFunction used to provide functions for Baryon Namespaces used inside
// roxygen2 return.
type ReturnFunction func(*tool.Outputs, string)

// returnInstruction is a map of functions used when parsing roxygen2 return.
var returnInstructions map[string]ReturnFunction = map[string]ReturnFunction{
	"data": func(o *tool.Outputs, args string) {
		argList := strings.Split(args, ",")
		if len(argList) < 2 {
			return
		}
		name := strings.TrimSpace(argList[0])
		format := strings.TrimSpace(argList[1])
		label := ""
		if len(argList) > 2 {
			label = strings.TrimSpace(argList[2])
		}
		newData := tool.Data{
			Format: format,
			Name:   name,
			Label:  label,
		}
		err := newData.Validate()
		if err != nil {
			log.Fatalf("parser_roxygen, returnInstruction[\"data\"]: %v", err)
			return
		}
		o.Data = append(o.Data, newData)
	},
}

var instructionRegex = regexp.MustCompile(`((?:[[:alpha:]]|!)+)\ *(?:\(([^)]*)|)`)

// Parse instruction into a tool.Param.
func parseInstruction(t *tool.Param, instruction string) {
	instructions := strings.Split(instruction, ";")
	for _, instruction := range instructions {
		instruction = strings.TrimSpace(instruction)
		match := instructionRegex.FindStringSubmatch(instruction)
		if len(match) < 3 {
			continue
		}
		parser, ok := paramIstructions[match[1]]
		if !ok {
			fmt.Fprintf(os.Stderr,
				"Instruction \"%s\" not found. Continuing.\n", match[1])
			continue
		}
		parser(t, match[2])
	}
}
