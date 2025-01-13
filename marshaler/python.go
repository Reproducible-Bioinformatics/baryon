package marshaler

import (
	"baryon/tool"
	"fmt"
	"strings"
)

// Ensure PythonMarshaler implements the Marshaler interface at compile-time.
var _ Marshaler = (*PythonMarshaler)(nil)

type PythonMarshaler struct{}

type PythonType struct {
	typeName  string
	typeCheck string
}

// Obtain a PythonType from a typeName of a tool.Param.
func (p PythonMarshaler) obtainType(typeName string, value string) (*PythonType, error) {
	switch typeName {
	case "text", "baseurl", "color", "file", "ftpfile", "hidden", "hidden_data":
		return &PythonType{
			typeName:  "str",
			typeCheck: fmt.Sprintf(`not isinstance(%s, str)`, value),
		}, nil
	case "integer":
		return &PythonType{
			typeName:  "int",
			typeCheck: fmt.Sprintf(`not isinstance(%s, int)`, value),
		}, nil
	case "float":
		return &PythonType{
			typeName:  "float",
			typeCheck: fmt.Sprintf(`not isinstance(%s, float)`, value),
		}, nil
	case "boolean":
		return &PythonType{
			typeName:  "bool",
			typeCheck: fmt.Sprintf(`%s not in [True, False]`, value),
		}, nil
	case "genomebuild", "select":
		return &PythonType{
			typeName:  "enum",
			typeCheck: fmt.Sprintf(`%s == ""`, value),
		}, nil
	case "data_column", "data", "data_collection", "drill_down":
		return &PythonType{
			typeName:  "file",
			typeCheck: fmt.Sprintf(`not os.path.isfile(%s)`, value),
		}, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}
}

// Marshal implements Marshaler.
func (p PythonMarshaler) Marshal(tool *tool.Tool) ([]byte, error) {
	buffer := []byte("")
	if out, err := p.marshalDescription(`"""%s"""`, tool.Description); err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	if out, err := p.marshalInputs(tool.Inputs); err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	if tool.Command == nil {
		return nil, fmt.Errorf("[PythonMarshaler.Marshal]: command not specified.")
	}
	if out, err := p.marshalContainerAndCommand(
		tool.Requirements.Container,
		*tool.Command,
	); err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.Marshal]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}

	return buffer, nil
}

func (p PythonMarshaler) marshalDescription(format string, description string) ([]byte, error) {
	buffer := []byte{}
	lines := strings.Split(description, "\n")
	for _, l := range lines {
		buffer = append(buffer, []byte(fmt.Sprintf(format, l))...)
		buffer = append(buffer, []byte("\n")...)
	}
	return buffer, nil
}

func (p PythonMarshaler) marshalInputs(inputs *tool.Inputs) ([]byte, error) {
	buffer := []byte("\n# Inputs\n")
	if out, err := p.processParams(inputs.Param); err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.marshalInputs]: %v", err)
	} else {
		buffer = append(buffer, out...)
	}
	return append(buffer, []byte("\n# End Inputs\n")...), nil
}

func (p PythonMarshaler) processParams(params []tool.Param) ([]byte, error) {
	if len(params) == 0 {
		return nil, nil
	}
	marshaledParam, err := p.marshalParam(&params[0])
	if err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.processParams]: %v", err)
	}
	remainingBytes, err := p.processParams(params[1:])
	if err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.processParams]: %v", err)
	}
	return append(marshaledParam, remainingBytes...), nil
}

func (p PythonMarshaler) marshalParam(param *tool.Param) ([]byte, error) {
	if param == nil {
		return nil, fmt.Errorf("[PythonMarshaler.marshalParam]: Empty field")
	}

	pythonType, err := p.obtainType(param.Type, param.Name)
	if err != nil {
		return nil, fmt.Errorf("[PythonMarshaler.marshalParam]: %v", err)
	}

	buffer := []byte(fmt.Sprintf("# %s\n", param.Help))
	buffer = append(buffer, []byte(fmt.Sprintf(
		`%s = None
if "--%s=" in args:
	%s = args["--%s="]
if %s:
	raise ValueError("%s is not of type %s")`,
		param.Name,
		param.Name,
		param.Name,
		param.Name,
		pythonType.typeCheck,
		param.Name,
		pythonType.typeName,
	))...)
	return buffer, nil
}

func (p PythonMarshaler) marshalContainerAndCommand(
	containers []tool.Container,
	command tool.Command,
) ([]byte, error) {
	buffer := []byte("# Command\n")
	for _, container := range containers {
		if container.Type != "docker" {
			return nil, fmt.Errorf("Only docker is supported")
		}
		buffer = append(buffer, []byte(
			fmt.Sprintf("subprocess.run(['docker', 'run', '--rm', %s'%s', f\"%s\"], check=True)\n",
				p.marshalVolumes(container.Volumes),
				container.Value,
				p.marshalCommand(command),
			))...)
	}
	return buffer, nil
}

func (p PythonMarshaler) marshalCommand(command tool.Command) string {
	buffer := []byte{}
	for _, element := range strings.Split(command.Value, " ") {
		if strings.HasPrefix(element, "$") {
			buffer = append(buffer, []byte(fmt.Sprintf(
				"{%s}", strings.TrimLeft(element, "$"),
			))...)
		} else {
			buffer = append(buffer, []byte(element)...)
		}
		buffer = append(buffer, []byte(" ")...)
	}
	return string(buffer)
}

func (p PythonMarshaler) marshalVolumes(mappings []tool.VolumeMapping) string {
	buffer := []byte{}
	for _, mapping := range mappings {
		hostPath := strings.TrimLeft(mapping.HostPath, "$")
		buffer = append(buffer, []byte(fmt.Sprintf(
			"f\" -v {%s}:%s\", ", hostPath, mapping.GuestPath,
		))...)
	}
	return string(buffer)
}
