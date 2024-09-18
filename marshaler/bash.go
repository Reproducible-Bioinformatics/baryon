package marshaler

import (
	"baryon/tool"
)

// Ensure bashMarshaler implements the Marshaler interface at compile-time.
var _ Marshaler = (*BashMarshaler)(nil)

type BashMarshaler struct{}

// Marshal implements Marshaler.
func (b BashMarshaler) Marshal(tool *tool.Tool) ([]byte, error) {
	buffer := []byte("#!/bin/bash\n")
	return buffer, nil
}
