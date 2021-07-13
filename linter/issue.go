package linter

// Issue is a structure of detected by linter problem.
type Issue struct {
	// File name is name of file, in which found issue.
	Filename string

	// LineNum - problem line number.
	LineNum int

	// Message is a text human-readable message of linter.
	Message string

	// LinterName - name of a linter that discovered the problem.
	LinterName string
}
