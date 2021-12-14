package main

import (
	"errors"
	"testing"

	"github.com/google/go-github/github"
)

func Test_ghClient_generateCommitBody(t *testing.T) {
	type args struct {
		pr *github.PullRequest
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate commit message with labels",
			args: args{
				pr: &github.PullRequest{
					Body: github.String("pull request body"),
					Labels: []*github.Label{
						{
							Name: github.String("label1"),
						},
						{
							Name: github.String("label2"),
						},
					},
				},
			},
			want: `
pull request body

Labels:
  * label1
  * label2`,
			wantErr: false,
		},
		{
			name: "generate commit message",
			args: args{
				pr: &github.PullRequest{
					Body: github.String("pull request body"),
				},
			},
			want: `
pull request body
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateCommitBody(tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ghClient.generateCommitMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ghClient.generateCommitMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCommitSubject(t *testing.T) {
	type args struct {
		pr *github.PullRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generate commit subject",
			args: args{
				pr: &github.PullRequest{
					Title:  github.String("pull request title"),
					Number: github.Int(1),
				},
			},
			want: "pull request title (#1)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateCommitSubject(tt.args.pr); got != tt.want {
				t.Errorf("generateCommitSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateEnv(t *testing.T) {
	type args struct {
		e env
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid env",
			args: args{
				e: env{
					Comment: "/merge",
					Mergers: []string{"0daryo"},
					Actor:   "0daryo",
				},
			},
		},
		{
			name: "invalid comment",
			args: args{
				e: env{
					Comment: "/approve",
					Mergers: []string{"0daryo"},
					Actor:   "0daryo",
				},
			},
			wantErr: true,
		},
		{
			name: "actor is not merger",
			args: args{
				e: env{
					Comment: "/merge",
					Mergers: []string{"0daryo"},
					Actor:   "github",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEnv(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("validateEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_errMsg(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{},
			want: "Succeeded!",
		},
		{
			name: "need approval",
			args: args{
				err: errors.New("failed to merge pull request: PUT https://api.github.com/repos/abema/github-actions-merger/pulls/1/merge: 405 At least 2 approving review is required by reviewers with write access. []"),
			},
			want: "Need 2 approving review",
		},
		{
			name: "internal server error",
			args: args{
				err: errors.New("internal server error"),
			},
			want: "internal server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errMsg(tt.args.err); got != tt.want {
				t.Errorf("errMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
