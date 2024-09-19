package marshaler

import (
	"baryon/tool"
	"fmt"
	"strings"
)

// Ensure bashMarshaler implements the Marshaler interface at compile-time.
var _ Marshaler = (*BashMarshaler)(nil)

type BashMarshaler struct{}

type BashType struct {
	typeName  string
	typeCheck string
}

// Obtain a BashType from a typeName of a tool.Param.
func (b BashMarshaler) obtainType(typeName string, value string) (*BashType, error) {
	switch typeName {
	case "text", "baseurl", "color", "file", "ftpfile", "hidden", "hidden_data":
		return &BashType{
			typeName:  "string",
			typeCheck: fmt.Sprintf("[[ ! -n $%s ]]", value),
		}, nil
	case "integer":
		return &BashType{
			typeName:  "int",
			typeCheck: fmt.Sprintf("[[ ! $%s =~ ^-?[0-9]+$ ]]", value),
		}, nil
	case "float":
		return &BashType{
			typeName:  "float",
			typeCheck: fmt.Sprintf("[[ ! $%s =~ ^-?[0-9]*\\.?[0-9]+$ ]]", value),
		}, nil
	case "boolean":
		return &BashType{
			typeName:  "bool",
			typeCheck: fmt.Sprintf("[[ ! $%s == \"true\" && ! $%s == \"false\" ]]", value, value),
		}, nil
	case "genomebuild", "select":
		return &BashType{
			typeName:  "enum",
			typeCheck: fmt.Sprintf("[[ $%s == \"\" ]]", value),
		}, nil
	case "data_column", "data", "data_collection", "drill_down":
		return &BashType{
			typeName:  "file",
			typeCheck: fmt.Sprintf("[[ ! -f $%s ]]", value),
		}, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}
}

// Marshal implements Marshaler.
func (b BashMarshaler) Marshal(tool *tool.Tool) ([]byte, error) {
	buffer := []byte{}
	echoDescription := []byte{}

	if out, err := b.marshalDescription(`echo "%s"
`, tool.Description); err != nil {
		return nil, fmt.Errorf("[bashMarshaler.Marshal]: %v", err)
	} else {
		echoDescription = out
	}
	buffer = append(buffer, []byte(fmt.Sprintf(`#!/bin/bash

usage() {
	%s
	echo "$0 usage:" && grep "\-.*)\ #" $0; exit 0;
}
[ $# -eq 0 ] && usage
`, echoDescription))...)

	if out, err := b.marshalDescription("# %s\n", tool.Description); err != nil {
		return nil, fmt.Errorf("[bashMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	if out, err := b.marshalInputs(tool.Inputs); err != nil {
		return nil, fmt.Errorf("[bashMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	return buffer, nil
}
func (b BashMarshaler) marshalDescription(description string) ([]byte, error) {
	buffer := []byte("\n")
	lines := strings.Split(description, "\n")
	for _, l := range lines {
		buffer = append(buffer, fmt.Sprintf("# %s\n", l)...)
	}
	return buffer, nil
}
	return buffer, nil
}
