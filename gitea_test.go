package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitea_SendReview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	gitea := Gitea{
		BaseURL: server.URL,
		Client:  server.Client(),
	}
	review := &Review{
		Body: "Found 100500 bugs",
		Comments: []ReviewComment{{
			Body:        "gitea.go:33:28: response body must be closed (bodyclose)",
			NewPosition: 33,
		}},
	}

	if err := gitea.SendReview("octocat/helloworld", 0, review); err != nil {
		t.Fatalf("%+v\n", err)
	}
}
