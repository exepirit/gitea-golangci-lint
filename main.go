package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/exepirit/gitea-golangci-lint/linter"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "gitea-golangci-lint",
	Usage: "Send linter outputs as pull request review",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "giteaUrl",
			Usage:   "Gitea server url",
			EnvVars: []string{"GITEA_URL", "PLUGIN_URL"},
		},
		&cli.StringFlag{
			Name:    "user",
			Usage:   "Gitea user name",
			EnvVars: []string{"GITEA_USER", "PLUGIN_USER"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "Gitea access token",
			EnvVars: []string{"GITEA_TOKEN", "PLUGIN_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   "Repository name, which is inspected. E. g. octocat/hello_world",
			EnvVars: []string{"GITEA_REPO", "DRONE_REPO"},
		},
		&cli.IntFlag{
			Name:    "pullRequest",
			Usage:   "Pull Request ID",
			EnvVars: []string{"PULL_REQUEST", "DRONE_PULL_REQUEST"},
		},
		&cli.IntFlag{
			Name:    "httpTimeout",
			Usage:   "HTTP request timeout in seconds",
			EnvVars: []string{"HTTP_TIMEOUT"},
			Value:   30,
		},
		&cli.StringFlag{
			Name:    "format",
			Usage:   "Input format of lint result, value will be empty or checkstyle, default to be empty",
			EnvVars: []string{"LINT_FORMAT"},
			Value:   "",
		},
	},
	HideVersion: true,
	Action:      lint,
}

func lint(ctx *cli.Context) error {
	repo, pullRequestID := ctx.String("repo"), ctx.Int("pullRequest")
	format := ctx.String("format")
	gitea := Gitea{
		BaseURL: strings.TrimSuffix(ctx.String("giteaUrl"), "/"),
		Client: &http.Client{
			Timeout: time.Duration(ctx.Int("httpTimeout")) * time.Second,
		},
		Username: ctx.String("user"),
		Token:    ctx.String("token"),
	}
	var issues []linter.Issue
	if format == "" {
		scanner := linter.NewLineScanner(os.Stdin)
		for scanner.Next() {
			issues = append(issues, scanner.Get())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scan linter output: %w", err)
		}
	} else if format == "checkstyle" {
		var err error
		issues, err = linter.NewCheckstyleScanner(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("format error: format value should be empty or checkstyle")
	}

	err := gitea.DiscardPreviousReviews(repo, pullRequestID)
	if err != nil {
		return fmt.Errorf("reset previous review: %w", err)
	}

	err = gitea.SendReview(repo, pullRequestID, FormatReview(issues))
	if err != nil {
		return fmt.Errorf("push new automated review: %w", err)
	}

	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
