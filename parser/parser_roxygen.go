package parser

import (
	"baryon/tool"
	"fmt"
	"log"
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
		if len(baryonInstruction) < 1 {
			return nil
		}
		for _, option := range strings.Split(baryonInstruction[1], ";") {
			option = strings.TrimSpace(option)
			match := instructionRegex.FindStringSubmatch(option)
			if len(match) < 3 {
				continue
			}
			option = strings.TrimSpace(match[1])
			if optionFunction, ok := paramOptions[option]; ok {
				optionFunction(&newParam, match[2])
			} else {
				return fmt.Errorf(`act["param"]: option "%s" not found.`, option)
			}
		}
		err := newParam.Validate()
		if err != nil {
			return fmt.Errorf(`act["param"]: %v`, err)
		}
		t.Inputs.Param = append(t.Inputs.Param, newParam)
		return nil
	},
	"description": func(description string, t *tool.Tool) error {
		cleanup, err := runInstruction(description, t, descriptionInstruction)
		if err != nil {
			return fmt.Errorf(`act["description"]: %v`, err)
		}
		t.Description = cleanup
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
		_, err := runInstruction(description, t, returnInstructions)
		if err != nil {
			return fmt.Errorf(`act["return"]: %v`, err)
		}
		return nil
	},
}

// runInstruction runs the parser and
func runInstruction(
	description string,
	t *tool.Tool,
	instruct map[string]ToolFunction,
) (string, error) {
	baryonInstruction :=
		baryonNamespaceRegex.FindStringSubmatch(description)
	if len(baryonInstruction) < 1 {
		return strings.TrimSpace(description), nil
	}
	// Processing of the description string according to Galaxy's specs.
	description =
		strings.Replace(description, baryonInstruction[0], "", -1)
	if err := parseInstruction(t, baryonInstruction[1], instruct); err != nil {
		return "", fmt.Errorf(`runInstruction: %v`, err)
	}
	return strings.TrimSpace(description), nil
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

// paramOptions is a map of function used to when parsing a roxygen2 param.
var paramOptions map[string]ParamFunction = map[string]ParamFunction{
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
var descriptionInstruction map[string]ToolFunction = map[string]ToolFunction{
	"container": func(t *tool.Tool, args string) error {
		argList := strings.Split(args, ",")
		if len(argList) < 1 {
			return fmt.Errorf("descriptionInstruction[\"container\"]: less than 1 arg")
		}
		container := tool.Container{
			Type:  "docker", // This is the default for baryon.
			Value: strings.TrimSpace(argList[0]),
		}
		if len(argList) > 1 {
			container.Type = strings.TrimSpace(argList[1])
		}
		if err := container.Validate(); err != nil {
			return fmt.Errorf("descriptionInstruction[\"container\"]: %v", err)
		}
		if t.Requirements == nil {
			t.Requirements = &tool.Requirements{}
		}
		t.Requirements.Container = append(t.Requirements.Container, container)
		return nil
	},
	"command": func(t *tool.Tool, arg string) error {
		arg = strings.TrimSpace(arg)
		if len(arg) == 0 {
			return fmt.Errorf(
				`descriptionInstruction["command"]: argument not present.`)
		}
		if t.Command == nil {
			t.Command = &tool.Command{Value: arg}
		}
		t.Command.Value = arg
		return nil
	},
	"volume": func(t *tool.Tool, args string) error {
		args = strings.TrimSpace(args)
		if len(args) == 0 {
			return fmt.Errorf(
				`descriptionInstruction["volume"]: argument not present.`)
		}
		argList := strings.Split(args, ":")
		if len(argList) != 2 {
			return fmt.Errorf(
				`descriptionInstruction["volume"]: volume should contain 2 arguments.`)
		}
		volMapping := tool.VolumeMapping{
			HostPath:  strings.TrimSpace(argList[0]),
			GuestPath: strings.TrimSpace(argList[1]),
		}
		if len(t.Requirements.Container) == 0 {
			return fmt.Errorf(
				`descriptionInstruction["volume"]: there's no container.`)
		}
		for i := range t.Requirements.Container {
			t.Requirements.Container[i].Volumes =
				append(t.Requirements.Container[i].Volumes, volMapping)
		}
		return nil
	},
}

// ToolFunction is used to provide functions for Baryon Namespaces used, for
// example, inside roxygen2 tags.
type ToolFunction func(t *tool.Tool, args string) error

// retrieveParser gets an instruction and instructions.
// If it doesn't find the instruction inside the map, it returns an error.
func retrieveParser(
	instruction string,
	instructions map[string]ToolFunction,
) (ToolFunction, error) {
	if instruction, ok := instructions[instruction]; ok {
		return instruction, nil
	}
	return nil,
		fmt.Errorf("evaluate: Instruction \"%s\" not found", instruction)
}

// returnInstruction is a map of functions used when parsing roxygen2 return.
var returnInstructions map[string]ToolFunction = map[string]ToolFunction{
	"data": func(o *tool.Tool, args string) error {
		argList := strings.Split(args, ",")
		if len(argList) < 2 {
			return fmt.Errorf("returnInstructions[\"data\"]: less than 2 args")
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
			return fmt.Errorf("returnInstructions[\"data\"]: %v", err)
		}
		if o.Outputs == nil {
			o.Outputs = &tool.Outputs{}
		}
		o.Outputs.Data = append(o.Outputs.Data, newData)
		return nil
	},
}

// instructionRegex is used to match a Baryon Instruction and obtain its name
// and the argument list.
var instructionRegex = regexp.MustCompile(`((?:[[:alpha:]]|!)+)\ *(?:\(([^)]*)|)`)

// Parse instruction into a tool.Tool.
func parseInstruction(
	t *tool.Tool,
	instructionList string,
	instructionMap map[string]ToolFunction,
) error {
	instructions := strings.Split(instructionList, ";")
	for _, instruction := range instructions {
		instruction = strings.TrimSpace(instruction)
		match := instructionRegex.FindStringSubmatch(instruction)
		if len(match) < 3 {
			continue
		}
		parser, err := retrieveParser(strings.TrimSpace(match[1]), instructionMap)
		if err != nil {
			return fmt.Errorf("parseInstruction: %v", err)
		}
		err = parser(t, match[2])
		if err != nil {
			return fmt.Errorf("parseInstruction: %v", err)
		}
	}
	return nil
}
