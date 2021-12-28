package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"text/template"
	"time"

	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

type env struct {
	GithubToken   string   `envconfig:"GITHUB_TOKEN"`
	Owner         string   `envconfig:"OWNER"`
	Repo          string   `envconfig:"REPO"`
	PRNumber      int      `envconfig:"PR_NUMBER"`
	Comment       string   `envconfig:"COMMENT"`
	MergeMethod   string   `envconfig:"MERGE_METHOD" default:"merge"`
	Mergers       []string `envconfig:"MERGERS"`
	AutoApprovers []string `envconfig:"AUTO_APPROVERS"`
	Actor         string   `envconfig:"GITHUB_ACTOR"` // github user who initiated the workflow.
}

const (
	mergeComment = "/merge"
	jobTimeout   = 10 * 60 * time.Second
)

func main() {
	var e env
	err := envconfig.Process("INPUT", &e)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := validateEnv(e); err != nil {
		log.Fatal(err.Error())
	}
	ctx, f := context.WithTimeout(context.Background(), jobTimeout)
	defer f()
	client := newGHClient(e.GithubToken)
	if autoApproveActor(e) {
		if err := client.approveIfBlocked(ctx, e.Owner, e.Repo, e.PRNumber); err != nil {
			if serr := client.sendMsg(ctx, e.Owner, e.Repo, e.PRNumber, errMsg(err)); serr != nil {
				log.Fatalf("failed to send message: %v original: %v", serr, err)
			}
			log.Fatal(err.Error())
		}
	}
	if err := client.merge(ctx, e.Owner, e.Repo, e.PRNumber, e.MergeMethod); err != nil {
		if serr := client.sendMsg(ctx, e.Owner, e.Repo, e.PRNumber, errMsg(err)); serr != nil {
			log.Fatalf("failed to send message: %v original: %v", serr, err)
		}
		log.Fatal(err.Error())
	}
	successMsg := "Merged PR #" + fmt.Sprintf("%d", e.PRNumber) + " successfully!"
	if err := client.sendMsg(ctx, e.Owner, e.Repo, e.PRNumber, successMsg); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf(successMsg)
}

func validateEnv(e env) error {
	if e.Comment != mergeComment {
		return fmt.Errorf("comment must be %s, got %s", mergeComment, e.Comment)
	}
	if len(e.Mergers) == 0 {
		return nil
	}
	for _, m := range e.Mergers {
		if e.Actor == m {
			return nil // if actor matches specified mergers, then valid workflow run
		}
	}
	return fmt.Errorf("actor %s is not in mergers list", e.Actor)
}

func autoApproveActor(e env) bool {
	if len(e.AutoApprovers) == 0 {
		log.Print("auto approvers list is empty")
		return false
	}
	log.Printf("try to check auto approvers %v", e.AutoApprovers)
	for _, aa := range e.AutoApprovers {
		if e.Actor == aa {
			log.Printf("actor %s is auto approver", e.Actor)
			return true // if actor matches specified auto approvers, returns true
		}
	}
	log.Printf("actor %s is not included in the auto approvers list", e.Actor)
	return false
}

type ghClient struct {
	client *github.Client
}

func newGHClient(token string) *ghClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &ghClient{
		client: client,
	}
}

func (gh *ghClient) approveIfBlocked(ctx context.Context, owner, repo string, prNumber int) error {
	log.Print("try to get pr for check to mergeable.")
	pr, _, err := gh.client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}
	if pr.GetMerged() {
		log.Print("this pr is already merged.")
		return nil
	}
	if pr.GetMergeableState() != "blocked" {
		log.Print("this pr mergeable state is not blocked.")
		return nil
	}
	log.Print("try to create review with approve.")
	event := "APPROVE"
	_, _, err = gh.client.PullRequests.CreateReview(ctx, owner, repo, prNumber, &github.PullRequestReviewRequest{Event: &event})
	if err != nil {
		return fmt.Errorf("failed to approve pull request: %w", err)
	}
	return nil
}

func (gh *ghClient) merge(ctx context.Context, owner, repo string, prNumber int, mergeMethod string) error {
	pr, _, err := gh.client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}
	commitMsg, err := generateCommitBody(pr)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}
	_, _, err = gh.client.PullRequests.Merge(ctx, owner, repo, prNumber, commitMsg, &github.PullRequestOptions{
		CommitTitle: generateCommitSubject(pr),
		MergeMethod: mergeMethod,
	})
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}
	return nil
}

func generateCommitSubject(pr *github.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", pr.GetTitle(), pr.GetNumber())
}

func generateCommitBody(pr *github.PullRequest) (string, error) {
	body := newCommitBody(pr)
	o := new(bytes.Buffer)
	if err := bodyTpl.Execute(o, body); err != nil {
		return "", err
	}
	return o.String(), nil
}

func (gh *ghClient) sendMsg(ctx context.Context, owner, repo string, prNumber int, msg string) error {
	_, _, err := gh.client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{
		Body: &msg,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func newCommitBody(pr *github.PullRequest) commitBody {
	labels := make([]string, 0, len(pr.Labels))
	for _, l := range pr.Labels {
		labels = append(labels, l.GetName())
	}
	return commitBody{
		Message: pr.GetBody(),
		Labels:  labels,
	}
}

type commitBody struct {
	Labels  []string
	Message string
}

var bodyTpl = template.Must(template.New("commit").Parse(`
{{- if .Message }}
{{ .Message }}
{{- end }}
{{if .Labels}}
Labels:
{{- range .Labels }}
  * {{ . }}
{{- end -}}
{{- end -}}
`))

var (
	needApproveRegexp = regexp.MustCompile("At least ([0-9]+) approving review is required by reviewers with write access")
)

// errMsg returns error message to post from error.
// Especially handing error from github. go-github does not have error type for some cases.
func errMsg(err error) string {
	if err == nil {
		return "Succeeded!"
	}
	ss := needApproveRegexp.FindStringSubmatch(err.Error())
	if len(ss) == 2 {
		return fmt.Sprintf("Need %s approving review", ss[1])
	}
	return err.Error()
}
