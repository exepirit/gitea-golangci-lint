package main

import (
	"io"
)

func ReadIssues(r io.Reader) []Issue {
	issues := make([]Issue, 0)

	issueScanner := NewLineScanner(r)
	for issueScanner.Next() {
		issues = append(issues, issueScanner.Get())
	}

	if err := issueScanner.Err(); err != nil {
		panic(err)
	}

	return issues
}
