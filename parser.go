package main

import (
	"regexp"
	"strconv"
)

var linterOutputRegex = regexp.MustCompile(`(.+?):(\d+):(\d*):? (.+) \((.+)\)`)

func ParseLinterOutput(output string) []Issue {
	matches := linterOutputRegex.FindAllStringSubmatch(output, -1)
	issues := make([]Issue, len(matches))
	for i, match := range matches {
		issues[i] = parseIssue(match)
	}
	return issues
}

func parseIssue(match []string) Issue {
	lineNum, err := strconv.Atoi(match[2])
	if err != nil {
		lineNum = 0
	}

	return Issue{
		Filename:   match[1],
		LineNum:    lineNum,
		Message:    match[4],
		LinterName: match[5],
	}
}
