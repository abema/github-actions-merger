package main

import (
	"errors"
	"testing"

	"github.com/google/go-github/github"
)

func Test_ghClient_generateCommitBody(t *testing.T) {
	type args struct {
		pr          *github.PullRequest
		gitTrailers []string
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
			want: `---
Labels:
  * label1
  * label2
---
pull request body
---
` +
				"```release-note\n" +
				"NONE\n" +
				"```",
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
---
pull request body
---
` +
				"```release-note\n" +
				"NONE\n" +
				"```",
			wantErr: false,
		},
		{
			name: "with release-note",
			args: args{
				pr: &github.PullRequest{
					Body: github.String("pull request body\n```release-note\nThis is greate a release!!!\n```"),
				},
			},
			want: `
---
pull request body

---
` +
				"```release-note\n" +
				"This is greate a release!!!\n" +
				"```",
			wantErr: false,
		},
		{
			name: "with git trailers",
			args: args{
				pr: &github.PullRequest{
					Body: github.String("pull request body"),
				},
				gitTrailers: []string{
					"Co-authored-by=abema",
					"Co-authored-by=actions",
				},
			},
			want: `Co-authored-by: abema
Co-authored-by: actions

---
pull request body
---
` +
				"```release-note\n" +
				"NONE\n" +
				"```",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateCommitBody(tt.args.pr, tt.args.gitTrailers)
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
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

func Test_splitReleaseNote(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name            string
		args            args
		wantDescription string
		wantReleaseNote string
	}{
		{
			name: "release note description",
			args: args{
				body: "release note description ```release-note\nThis is great release!!!\n```",
			},
			wantDescription: "release note description ",
			wantReleaseNote: "This is great release!!!",
		},
		{
			name: "no releaes note",
			args: args{
				body: "release note description",
			},
			wantDescription: "release note description",
			wantReleaseNote: "NONE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDescription, gotReleaseNote := splitReleaseNote(tt.args.body)
			if gotDescription != tt.wantDescription {
				t.Errorf("splitReleaseNote() gotDescription = %v, want %v", gotDescription, tt.wantDescription)
			}
			if gotReleaseNote != tt.wantReleaseNote {
				t.Errorf("splitReleaseNote() gotReleaseNote = %v, want %v", gotReleaseNote, tt.wantReleaseNote)
			}
		})
	}
}
