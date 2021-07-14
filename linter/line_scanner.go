package linter

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
)

var lineOutputIssueRegex = regexp.MustCompile(`(.+?):(\d+):(\d*):? (.+) \((.+)\)`)

func NewLineScanner(r io.Reader) IssueScanner {
	return &issueLineScanner{
		reader: bufio.NewReader(r),
	}
}

type issueLineScanner struct {
	reader  *bufio.Reader
	err     error
	scanned []Issue
}

func (sc *issueLineScanner) Next() bool {
	for len(sc.scanned) == 0 {
		line, err := sc.reader.ReadString('\n')
		switch {
		case len(line) == 0 || err == io.EOF:
			return false
		case err != nil:
			sc.catchError(err)
			return false
		}
		sc.scanned = append(sc.scanned, sc.parseLine(line)...)
	}
	return true
}

func (issueLineScanner) parseLine(line string) []Issue {
	matches := lineOutputIssueRegex.FindAllStringSubmatch(line, 1)
	issues := make([]Issue, len(matches))
	for i, match := range matches {
		lineNum, err := strconv.Atoi(match[2])
		if err != nil {
			lineNum = 0
		}

		issues[i] = Issue{
			Filename:   match[1],
			LineNum:    lineNum,
			Message:    match[4],
			LinterName: match[5],
		}
	}
	return issues
}

func (sc *issueLineScanner) catchError(err error) {
	sc.err = err
}

func (sc *issueLineScanner) Err() error {
	return sc.err
}

func (sc *issueLineScanner) Get() Issue {
	if len(sc.scanned) < 1 {
		panic("scanner does not contain any issue")
	}
	issue := sc.scanned[0]
	sc.scanned = sc.scanned[1:]
	return issue
}
