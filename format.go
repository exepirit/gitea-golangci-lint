package main

import (
	"fmt"

	"github.com/exepirit/gitea-golangci-lint/linter"
)

const reviewBodyTemplate = `Linter found %d issues.`
const commentBodyTemplate = `🔔 __%s__ found issue at line %d:
> %s`

func FormatReview(issues []linter.Issue) *Review {
	if len(issues) > 0 {
		return formatRequestChangesReview(issues)
	}
	return formatApproveReview()
}

func formatRequestChangesReview(issues []linter.Issue) *Review {
	comments := make([]ReviewComment, len(issues))
	for i, issue := range issues {
		body := fmt.Sprintf(commentBodyTemplate, issue.LinterName, issue.LineNum, issue.Message)
		comments[i] = ReviewComment{
			Body:        body,
			NewPosition: issue.LineNum,
			Path:        issue.Filename,
		}
	}
	return &Review{
		Body:     fmt.Sprintf(reviewBodyTemplate, len(issues)),
		Comments: comments,
		Event:    ReviewStateRequestChanges,
	}
}

func formatApproveReview() *Review {
	return &Review{
		Body:     fmt.Sprintf(reviewBodyTemplate, 0),
		Comments: []ReviewComment{},
		Event:    ReviewStateApproved,
	}
}
