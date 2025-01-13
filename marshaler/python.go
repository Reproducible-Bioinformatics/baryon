package marshaler

import "baryon/tool"

// Ensure bashMarshaler implements the Marshaler interface at compile-time.
var _ Marshaler = (*PythonMarshaler)(nil)

type PythonMarshaler struct{}

// Marshal implements Marshaler.
func (p PythonMarshaler) Marshal(*tool.Tool) ([]byte, error) {
	panic("unimplemented")
}
