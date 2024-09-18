package marshaler

import "baryon/tool"

// Marshaler defines an interface for serializing a tool.Tool instance.
// Implementations of this interface are responsible for taking a
// tool.Tool object and transforming it into a byte slice ([]byte).
type Marshaler interface {
	// Marshal converts a tool.Tool instance into a serialized byte slice.
	// Returns the byte representation of the tool or an error if the
	// marshalling process fails.
	//
	// Params:
	//    tool.Tool: The tool instance to be serialized.
	//
	// Returns:
	//    []byte: The serialized byte slice of the tool.Tool instance.
	//    error: An error object in case of a failure during marshalling.
	Marshal(*tool.Tool) ([]byte, error)
}
