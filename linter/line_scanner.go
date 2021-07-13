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
	if len(sc.scanned) > 0 {
		return true
	}
	sc.scan()
	return len(sc.scanned) > 0 && sc.err == nil
}

func (sc *issueLineScanner) scan() {
	line, err := sc.reader.ReadString('\n')
	if err != nil {
		sc.err = err
		if err != io.EOF {
			return
		}
	}

	parsed := sc.parseLine(line)
	sc.scanned = append(sc.scanned, parsed...)
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

func (sc *issueLineScanner) Err() error {
	if sc.err == io.EOF {
		return nil
	}
	return sc.err
}

func (sc *issueLineScanner) Get() Issue {
	if len(sc.scanned) == 0 {
		panic("scanner does not contain any issue")
	}
	issue := sc.scanned[0]
	sc.scanned = sc.scanned[1:]
	return issue
}
