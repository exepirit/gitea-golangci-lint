package linter

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// DefaultCheckStyleVersion defines the default "version" attribute on "<checkstyle>" lememnt
var DefaultCheckStyleVersion = "1.0.0"

// Severity defines a checkstyle severity code
type Severity string

var (
	SeverityError   Severity = "error"
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityIgnore  Severity = "ignore"
	SeverityNone    Severity
)

// CheckStyle represents a <checkstyle> xml element found in a checkstyle_report.xml file.
type CheckStyle struct {
	XMLName xml.Name `xml:"checkstyle"`
	Version string   `xml:"version,attr"`
	File    []*File  `xml:"file"`
}

// AddFile adds a checkstyle.File with the given filename.
func (cs *CheckStyle) AddFile(csf *File) {
	cs.File = append(cs.File, csf)
}

// GetFile gets a CheckStyleFile with the given filename.
func (cs *CheckStyle) GetFile(filename string) (csf *File, ok bool) {
	for _, file := range cs.File {
		if file.Name == filename {
			csf = file
			ok = true
			return
		}
	}
	return
}

// EnsureFile ensures that a CheckStyleFile with the given name exists
// Returns either an exiting CheckStyleFile (if a file with that name exists)
// or a new CheckStyleFile (if a file with that name does not exists)
func (cs *CheckStyle) EnsureFile(filename string) (csf *File) {
	csf, ok := cs.GetFile(filename)
	if !ok {
		csf = NewFile(filename)
		cs.AddFile(csf)
	}
	return csf
}

// String implements Stringer. Returns as xml.
func (cs *CheckStyle) String() string {
	checkStyleXML, err := xml.Marshal(cs)
	if err != nil {
		panic(err)
	}
	return string(checkStyleXML)
}

// New returns a new CheckStyle
func New() *CheckStyle {
	return &CheckStyle{Version: DefaultCheckStyleVersion, File: []*File{}}
}

// File represents a <file> xml element.
type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Error   []*Error `xml:"error"`
}

// AddError adds a checkstyle.Error to the file.
func (csf *File) AddError(cse *Error) {
	csf.Error = append(csf.Error, cse)
}

// NewFile creates a new checkstyle.File
func NewFile(filename string) *File {
	return &File{Name: filename, Error: []*Error{}}
}

// Error represents a <error> xml element
type Error struct {
	XMLName  xml.Name `xml:"error"`
	Line     int      `xml:"line,attr"`
	Column   int      `xml:"column,attr,omitempty"`
	Severity Severity `xml:"severity,attr,omitempty"`
	Message  string   `xml:"message,attr"`
	Source   string   `xml:"source,attr"`
}

// NewError creates a new checkstyle.Error
// Note that line starts at 0, and column starts at 1
func NewError(line int, column int, severity Severity, message string, source string) *Error {
	return &Error{Line: line, Column: column, Severity: severity, Message: message, Source: source}
}

// ParseCheckStyleFromXml Parse checkstyle xml to gitea issues.
func ParseCheckStyleFromXml(checkStyleXML []byte) ([]Issue, error) {
	issues := make([]Issue, 0)
	if len(checkStyleXML) == 0 {
		return issues, nil
	}
	checkStyle := New()
	err := xml.Unmarshal(checkStyleXML, checkStyle)
	if err == nil {
		for _, file := range checkStyle.File {
			for _, codingError := range file.Error {
				issues = append(issues, Issue{
					Filename:   file.Name,
					LineNum:    codingError.Line,
					Message:    codingError.Message,
					LinterName: codingError.Source,
				})
			}
		}
	}
	return issues, err
}

// NewCheckstyleScanner Create issues from io.Reader, eg: os.Stdin
func NewCheckstyleScanner(r io.Reader) ([]Issue, error) {
	readText, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %s", err)
	}
	return ParseCheckStyleFromXml(readText)
}
