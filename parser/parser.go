package parser

import "baryon/tool"

type Parser interface {
	// Parse parses a []byte.
	Parse([]byte) (*tool.Tool, error)
}

// rlangDocParser implements the functions to parse R function documentation
// and obtain a Galaxy Tool.
type rlangDocParser struct{}

// NewRlangDocParser returns a New rlangDocParser.
func NewRlangDocParser() *rlangDocParser {
	return &rlangDocParser{}
}

func (*rlangDocParser) Parse(in []byte) (*tool.Tool, error) {
	return &tool.Tool{}, nil
}
