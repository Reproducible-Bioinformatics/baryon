package tool

import "encoding/xml"

type Tool struct {
	XMLName     xml.Name    `xml:"tool"`
	Description Description `xml:"description"`
}

type Description string
