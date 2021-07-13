package main

import (
	"io"

	"github.com/exepirit/gitea-golangci-lint/linter"
)

func ReadIssues(r io.Reader) []linter.Issue {
	issues := make([]linter.Issue, 0)

	issueScanner := linter.NewLineScanner(r)
	for issueScanner.Next() {
		issues = append(issues, issueScanner.Get())
	}

	if err := issueScanner.Err(); err != nil {
		panic(err)
	}

	return issues
}
