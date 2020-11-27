package main

import (
	"bufio"
	"io"
	"os"
)

func ReadIssues(r io.Reader) []Issue {
	inputReader := bufio.NewReader(r)
	issues := make([]Issue, 0)
	for {
		line, err := inputReader.ReadBytes('\n')
		if err != nil && err == io.EOF {
			break
		}
		_, _ = os.Stdout.Write(line)

		newIssues := ParseLinterOutput(string(line))
		if len(newIssues) > 0 {
			issues = append(issues, newIssues...)
		}
	}
	return issues
}
