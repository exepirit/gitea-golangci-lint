package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const (
	StatusSuccess = "success"
	StatusFailure = "failure"
)

type Gitea struct {
	BaseURL  string
	Client   *http.Client
	Username string
	Token    string
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("server returns status %d", resp.StatusCode)
	}
	return nil
}

func (gitea Gitea) putHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func (gitea Gitea) SendReview(repo string, pullIndex int, review *Review) error {
	reqBody, err := json.Marshal(review)
	if err != nil {
		return errors.WithStack(err)
	}

	url := fmt.Sprintf("%s/api/v1/repos/%s/pulls/%d/reviews?access_token=%s", gitea.BaseURL, repo, pullIndex, gitea.Token)
	log.Println("Send Review To: " + url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.WithStack(err)
	}

	gitea.putHeaders(req)
	resp, err := gitea.Client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return checkStatusCode(resp)
}

func (gitea Gitea) DiscardPreviousReviews(repo string, pullIndex int) error {
	reviews, err := gitea.GetAllReviews(repo, pullIndex)
	if err != nil {
		return err
	}

	for _, review := range reviews {
		if review.User["username"].(string) != gitea.Username {
			continue
		}
		if err := gitea.DiscardReview(repo, pullIndex, review.ID); err != nil {
			return err
		}
	}

	return nil
}

func (gitea Gitea) GetAllReviews(repo string, pullIndex int) ([]PullReview, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/pulls/%d/reviews?access_token=%s", gitea.BaseURL, repo, pullIndex, gitea.Token)
	log.Println("Get All Reviews: " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gitea.putHeaders(req)
	resp, err := gitea.Client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if err := checkStatusCode(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return make([]PullReview, 0), nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var reviews []PullReview
	if err := json.Unmarshal(body, &reviews); err != nil {
		return nil, errors.WithStack(err)
	}

	return reviews, nil
}

func (gitea Gitea) DiscardReview(repo string, pullIndex int, reviewIndex int) error {
	url := fmt.Sprintf("%s/api/v1/repos/%s/pulls/%d/reviews/%d?access_token=%s", gitea.BaseURL, repo, pullIndex, reviewIndex, gitea.Token)
	log.Println("Discard Review: " + url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	gitea.putHeaders(req)
	resp, err := gitea.Client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return checkStatusCode(resp)
}

func (gitea Gitea) GetPullRequest(repo string, pullIndex int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/pulls/%d?access_token=%s", gitea.BaseURL, repo, pullIndex, gitea.Token)
	log.Println("Get PullRequest: " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gitea.putHeaders(req)
	resp, err := gitea.Client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	if err := checkStatusCode(resp); err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pr := &PullRequest{}
	if err := json.Unmarshal(body, pr); err != nil {
		return nil, errors.WithStack(err)
	}

	return pr, nil
}

func (gitea Gitea) SendCommitStatus(repo string, pullIndex int, ctx, status string) error {
	pr, err := gitea.GetPullRequest(repo, pullIndex)
	if err != nil {
		return errors.WithStack(err)
	}
	stat := &CommitStatus{
		Context:     ctx,
		Description: fmt.Sprintf("checked this commit with %s result", status),
		State:       status,
		TargetURL:   fmt.Sprintf("%s/%s/pulls/%d", gitea.BaseURL, repo, pullIndex),
	}
	reqBody, err := json.Marshal(stat)
	if err != nil {
		return errors.WithStack(err)
	}

	url := fmt.Sprintf("%s/api/v1/repos/%s/statuses/%s?access_token=%s", gitea.BaseURL, repo, pr.Head.Sha, gitea.Token)
	log.Println("Send Commit Status To: " + url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.WithStack(err)
	}

	gitea.putHeaders(req)
	resp, err := gitea.Client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return checkStatusCode(resp)
}
