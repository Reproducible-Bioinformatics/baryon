package tool

import "encoding/xml"

// Tool provides a representation of a Galaxy Tool xml file schema.
//
// You can find the current schema here:
// https://docs.galaxyproject.org/en/master/dev/schema.html
type Tool struct {
	XMLName     xml.Name    `xml:"tool"`
	Description Description `xml:"description"`
}

type Description string

type EdamTopics struct {
	XMLName   xml.Name    `xml:"edam_topics"`
	EdamTopic []EdamTopic `xml:"edam_topic"`
}

type EdamTopic string

type EdamOperations struct {
	XMLName       xml.Name        `xml:"edam_operations"`
	EdamOperation []EdamOperation `xml:"edam_operation"`
}

type EdamOperation string

type Xrefs struct {
	XMLName xml.Name `xml:"xrefs"`
	Xref    []Xref   `xml:"xref"`
}

type Xref struct {
	XMLName xml.Name `xml:"xref"`
	Type    string   `xml:"type,attr"`
}

type Creator struct {
	XMLName xml.Name `xml:"creator"`
	Person  Person   `xml:"person"`
}

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	// TODO: Add other fields.
}

// TODO: Add organization.
