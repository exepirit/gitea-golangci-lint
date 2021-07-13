package linter

type IssueScanner interface {
	Next() bool
	Err() error
	Get() Issue
}
