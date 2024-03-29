package main

type ReviewStateType string

const (
	ReviewStateApproved       ReviewStateType = "APPROVED"
	ReviewStateRequestChanges ReviewStateType = "REQUEST_CHANGES"
)

type ReviewComment struct {
	Body        string `json:"body"`
	NewPosition int    `json:"new_position"`
	Path        string `json:"path"`
}

type Review struct {
	Body     string          `json:"body"`
	Comments []ReviewComment `json:"comments"`
	Event    ReviewStateType `json:"event"`
}

type PullReview struct {
	Body           string                 `json:"body"`
	CommentsCount  int                    `json:"comments_count"`
	CommitID       string                 `json:"commit_id"`
	HTMLURL        string                 `json:"html_url"`
	ID             int                    `json:"id"`
	Official       bool                   `json:"official"`
	PullRequestURL string                 `json:"pull_request_url"`
	Stale          bool                   `json:"stale"`
	State          string                 `json:"state"`
	SubmittedAt    string                 `json:"submitted_at"`
	User           map[string]interface{} `json:"user"`
}

type PullRequest struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Number    int    `json:"number"`
	Head      Head   `json:"head"`
	MergeBase string `json:"merge_base"`
}
type Head struct {
	Sha string `json:"sha"`
}

type CommitStatus struct {
	Context     string `json:"context"`
	Description string `json:"description"`
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
}
