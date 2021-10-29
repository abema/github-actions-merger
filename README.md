# github-actions-merger
github-actions-merger is github actions that merges pull request with commit message including pull request labels.

# labelcommit
labelcommit is github actions that merges pull request with commit message including pull request labels.

## Usage
1. Write your workflow file.
```
  - name: merge
    uses: abema/github-actions-merger@main
    with: 
      "github token": ${{ secrets.GITHUB_TOKEN }}
      "owner": ${{ github.event.repository.owner.login }}
      "repo": ${{ github.event.repository.name }}
      "pr number": ${{ github.event.issue.number }}
      "comment": ${{ github.event.comment.body }}
```
https://github.com/abema/github-actions-merger/blob/main/.github/workflows/github-actions-merger.yaml

2. comment ```/merge``` on github pull request comment.

3. pull request is merged, and commit message includes labels.
```
fix: readme
* documentation
* enhancement
```

## Parameters
You need to set parameters in workflow.
```
github: token: ${{ secrets.GITHUB_TOKEN }}
owner: repository owner
repo: repository name
pr number: ${{ github.event.issue.number }}
comment: ${{ github.event.comment.body }}
```
