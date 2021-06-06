package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
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
