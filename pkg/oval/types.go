package oval

import "encoding/xml"

// OVALDefinition represents the root OVAL element
type OVALDefinition struct {
	XMLName     xml.Name         `xml:"oval_definitions"`
	Definitions []Definition     `xml:"definitions>definition"`
	Tests       []interface{}    `xml:"tests"`
	Objects     []interface{}    `xml:"objects"`
	States      []interface{}    `xml:"states"`
}

// Definition represents a single OVAL definition
type Definition struct {
	ID       string        `xml:"id,attr"`
	Version  string        `xml:"version,attr"`
	Class    string        `xml:"class,attr"`
	Metadata Metadata      `xml:"metadata"`
	Criteria Criteria      `xml:"criteria"`
}

// Metadata contains definition metadata
type Metadata struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

// Criteria contains test references
type Criteria struct {
	Criterion []Criterion `xml:"criterion"`
}

// Criterion references a test
type Criterion struct {
	TestRef string `xml:"test_ref,attr"`
}

// RegistryObject represents a Windows registry object
type RegistryObject struct {
	ID    string `xml:"id,attr"`
	Hive  string `xml:"hive"`
	Key   string `xml:"key"`
	Name  string `xml:"name"`
}

// FileObject represents a file object
type FileObject struct {
	ID       string `xml:"id,attr"`
	Path     string `xml:"path"`
	Filename string `xml:"filename"`
}

// RegistryState represents expected registry state
type RegistryState struct {
	ID        string `xml:"id,attr"`
	Value     string `xml:"value"`
	Operation string `xml:"operation,attr"`
}

// TestResult represents the outcome of a test
type TestResult struct {
	TestID    string      `json:"test_id"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Timestamp string      `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}