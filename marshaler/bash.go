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

// Marshal implements Marshaler.
func (b BashMarshaler) Marshal(tool *tool.Tool) ([]byte, error) {
	buffer := []byte("#!/bin/bash\n")

	if out, err := b.marshalDescription(tool.Description); err != nil {
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
