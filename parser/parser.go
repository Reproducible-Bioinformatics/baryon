package parser

import (
	"baryon/tool"
)

type Parser interface {
	// Parse parses a []byte.
	Parse([]byte) (*tool.Tool, error)
}
