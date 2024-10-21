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
	if tool.Command == nil {
		return nil, fmt.Errorf("[bashMarshaler.Marshal]: command not specified.")
	}
	if out, err := b.marshalContainerAndCommand(
		tool.Requirements.Container,
		*tool.Command,
	); err != nil {
		return nil, fmt.Errorf("[bashMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	return buffer, nil
}

func (b BashMarshaler) marshalDescription(format string, description string) ([]byte, error) {
	buffer := []byte("\n")
	lines := strings.Split(description, "\n")
	for _, l := range lines {
		buffer = append(buffer, fmt.Sprintf(format, l)...)
	}
	return buffer, nil
}

func (b BashMarshaler) marshalInputs(inputs *tool.Inputs) ([]byte, error) {
	buffer := []byte("\n# Inputs\n")
	if out, err := b.processParams(inputs.Param); err != nil {
		return nil, fmt.Errorf("[bashMarshaler.marshalInputs]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}
	return append(buffer, []byte("\n# End Inputs\n")...), nil
}

func (b BashMarshaler) processParams(params []tool.Param) ([]byte, error) {
	if len(params) == 0 {
		return nil, nil
	}
	marshaledParam, err := b.marshalParam(&params[0])
	if err != nil {
		return nil, fmt.Errorf("[bashMarshaler.processParams]: %v", err)
	}
	remainingBytes, err := b.processParams(params[1:])
	if err != nil {
		return nil, fmt.Errorf("[bashMarshaler.processParams]: %v", err)
	}
	return append(marshaledParam, remainingBytes...), nil
}

func (b BashMarshaler) marshalParam(param *tool.Param) ([]byte, error) {
	if param == nil {
		return nil, fmt.Errorf("[bashMarshaler.marshalParam]: Empty field")
	}

	bashType, err := b.obtainType(param.Type, param.Name)
	if err != nil {
		return nil, fmt.Errorf("[BashMarshaler.marshalParam]: %v", err)
	}
	buffer := []byte(fmt.Sprintf("## %s", param.Name))

	// Obtain from a parameter
	buffer = append(buffer, []byte(fmt.Sprintf(`
%s=""

for arg in "$@"; do
	case $arg in
		--%s=*) # %s
		%s="${arg#*=}"
		shift
		;;
	esac
done

if %s; then
	echo "%s is not of type %s"
	%s
fi
`,
		param.Name,
		param.Name,
		param.Help,
		param.Name,
		bashType.typeCheck,
		param.Name,
		bashType.typeName,
		func() string {
			if param.Optional {
				return fmt.Sprintf(`echo "WARN: %s is optional"`, param.Name)
			}
			return fmt.Sprintf(`exit 1`)
		}(),
	))...)
	return buffer, nil
}

func (b BashMarshaler) marshalContainerAndCommand(
	containers []tool.Container,
	command tool.Command,
) ([]byte, error) {
	buffer := []byte("# Command\n")
	for _, container := range containers {
		if container.Type != "docker" {
			return nil, fmt.Errorf("Only docker is supported")
		}
		buffer = append(buffer, []byte(
			fmt.Sprintf("docker run %s --rm %s '%s'\n",
				b.marshalVolumes(container.Volumes),
				container.Value,
				command.Value,
			))...)
	}
	return buffer, nil
}

func (b BashMarshaler) marshalVolumes(mappings []tool.VolumeMapping) string {
	buffer := []byte{}
	for _, mapping := range mappings {
		buffer = append(buffer, []byte(fmt.Sprintf(
			"-v %s:%s", mapping.HostPath, mapping.GuestPath,
		))...)
	}
	return string(buffer)
}
