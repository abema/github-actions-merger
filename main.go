package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

type env struct {
	GithubToken string `envconfig:"GITHUB_TOKEN"`
	Owner       string `envconfig:"OWNER"`
	Repo        string `envconfig:"REPO"`
	PRNumber    int    `envconfig:"PR_NUMBER"`
	Comment     string `envconfig:"COMMENT"`
	MergeMethod string `envconfig:"MERGE_METHOD" default:"merge"`
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

	if e.Comment != mergeComment {
		log.Println("comment is not /merge")
		return
	}
	ctx, f := context.WithTimeout(context.Background(), jobTimeout)
	defer f()
	client := newGHClient(e.GithubToken)
	if err := client.merge(ctx, e.Owner, e.Repo, e.PRNumber, e.MergeMethod); err != nil {
		if serr := client.sendMsg(ctx, e.Owner, e.Repo, e.PRNumber, err.Error()); serr != nil {
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

func generateCommitMessage(labels []string) string {
	var commitMessage string
	if len(labels) > 0 {
		commitMessage = "* " + strings.Join(labels, "\n* ")
	}
	return commitMessage
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

func (gh *ghClient) merge(ctx context.Context, owner, repo string, prNumber int, mergeMethod string) error {
	pr, _, err := gh.client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}
	labels := make([]string, 0, len(pr.Labels))
	for _, l := range pr.Labels {
		labels = append(labels, l.GetName())
	}
	commitMsg := generateCommitMessage(labels)

	_, _, err = gh.client.PullRequests.Merge(ctx, owner, repo, prNumber, commitMsg, &github.PullRequestOptions{
		CommitTitle: pr.GetTitle(),
		MergeMethod: mergeMethod,
	})
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}
	return nil
}

func (gh *ghClient) sendMsg(ctx context.Context, owner, repo string, prNumber int, msg string) error {
	_, ghResp, err := gh.client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{
		Body: &msg,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w, githubResponse: %s", err, ghResp.String())
	}
	return nil
}
