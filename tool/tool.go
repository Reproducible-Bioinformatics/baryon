package tool

import "encoding/xml"

// Tool provides a representation of a Galaxy Tool xml file schema.
//
// You can find the current schema here:
// https://docs.galaxyproject.org/en/master/dev/schema.html
type Tool struct {
	XMLName xml.Name `xml:"tool"`
	// The value is displayed in the tool menu immediately following the hyperlink
	// for the tool (based on the name attribute of the <tool> tag set described
	// above).
	//
	// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-description
	Description    string          `xml:"description"`
	EdamTopics     *EdamTopics     `xml:"edam_topics,omitempty"`
	EdamOperations *EdamOperations `xml:"edam_operations,omitempty"`
	Xrefs          *Xrefs          `xml:"xrefs,omitempty"`
	Creator        *Creator        `xml:"creator,omitempty"`
	Requirements   *Requirements   `xml:"requirements"`
	Command        *Command        `xml:"command"`
	Inputs         *Inputs         `xml:"inputs"`
}

// Container tag set for the <edam_topic> tags. A tool can have any number of
// EDAM topic references.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-edam-topics
type EdamTopics struct {
	XMLName   xml.Name    `xml:"edam_topics,omitempty"`
	EdamTopic []EdamTopic `xml:"edam_topic,omitempty"`
}

type EdamTopic string

// Container tag set for the <edam_operation> tags. A tool can have any number
// of EDAM operation references.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-edam-operations
type EdamOperations struct {
	XMLName       xml.Name        `xml:"edam_operations"`
	EdamOperation []EdamOperation `xml:"edam_operation"`
}

type EdamOperation string

// Container tag set for the <xref> tags. A tool can refer multiple reference
// IDs.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-xrefs
type Xrefs struct {
	XMLName xml.Name `xml:"xrefs"`
	Xref    []Xref   `xml:"xref"`
}

// The xref element specifies reference information according to a catalog.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-xrefs-xref
type Xref struct {
	XMLName xml.Name `xml:"xref"`
	// Type of reference - currently bio.tools, bioconductor, and biii
	// are the only supported options.
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// The creator(s) of this work. See schema.org/creator.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-creator
type Creator struct {
	XMLName      xml.Name      `xml:"creator,omitempty"`
	Person       []Person      `xml:"person,omitempty"`
	Organization *Organization `xml:"organization,omitempty"`
}

// Describes a person. Tries to stay close to schema.org/Person.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-creator-person
type Person struct {
	XMLName xml.Name `xml:"person,omitempty"`
	Name    string   `xml:"name,omitempty"`
}

// Describes an organization. Tries to stay close to schema.org/Organization.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-creator-organization
type Organization struct {
	XMLName xml.Name `xml:"organization,omitempty"`
	Name    string   `xml:"name,omitempty"`
}

// This is a container tag set for the requirement, resource and container tags
// described in greater detail below. requirements describe software packages
// and other individual computing requirements required to execute a tool,
// while containers describe Docker or Singularity containers that should be
// able to serve as complete descriptions of the runtime of a tool.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-requirements
type Requirements struct {
	XMLName     xml.Name      `xml:"requirements"`
	Requirement []Requirement `xml:"requirement,omitempty"`
	Container   *Container    `xml:"container,omitempty"`
}

// This tag set is contained within the <requirements> tag set. Third party
// programs or modules that the tool depends upon are included in this tag set.
//
// When a tool runs, Galaxy attempts to resolve these requirements (also called
// dependencies). requirements are meant to be abstract and resolvable by
// multiple different dependency resolvers (e.g. conda, the Galaxy Tool Shed
// dependency management system, or environment modules).
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-requirements-requirement
type Requirement struct {
	XMLName xml.Name `xml:"requirement"`
	Type    string   `xml:"type,attr"`
	Version string   `xml:"version,attr"`
}

// This tag set is contained within the ‘requirements’ tag set. Galaxy can be
// configured to run tools within Docker or Singularity containers - this tag
// allows the tool to suggest possible valid containers for this tool.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-requirements-container
type Container struct {
	XMLName xml.Name `xml:"container"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}

// This tag specifies how Galaxy should invoke the tool’s executable, passing
// its required input parameter values (the command line specification links
// the parameters supplied in the form with the actual tool executable).
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-command
type Command struct {
	XMLName xml.Name `xml:"command"`
	Value   string   `xml:",cdata"`
}

// Consists of all elements that define the tool’s input parameters.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-inputs
type Inputs struct {
	XMLName xml.Name `xml:"inputs"`
	Param   []Param  `xml:"param"`
}

// Contained within the <inputs> tag set - each of these specifies a field that
// will be displayed on the tool form. Ultimately, the values of these form
// fields will be passed as the command line parameters to the tool’s
// executable.
//
// https://docs.galaxyproject.org/en/latest/dev/schema.html#tool-inputs-param
type Param struct {
	XMLName         xml.Name `xml:"param"`
	Type            string   `xml:"type"`
	Name            string   `xml:"name"`
	Value           string   `xml:"value"`
	Argument        string   `xml:"argument"`
	Label           string   `xml:"label"`
	Help            string   `xml:"help"`
	Optional        bool     `xml:"optional"`
	RefreshOnChange bool     `xml:"refresh_on_change"`
}
