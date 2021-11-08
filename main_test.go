package main

import (
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
