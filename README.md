# github-actions-merger
github-actions-merger is github actions that merges pull request with commit message including pull request labels.

## Usage
1. Write your workflow file.
```
  - name: merge
    uses: abema/github-actions-merger@main
    with: 
      "github_token": ${{ secrets.GITHUB_TOKEN }}
      "owner": ${{ github.event.repository.owner.login }}
      "repo": ${{ github.event.repository.name }}
      "pr_number": ${{ github.event.issue.number }}
      "comment": ${{ github.event.comment.body }}
      "mergers": 'na-ga,0daryo'
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
github_token: ${{ secrets.GITHUB_TOKEN }}
owner: ${{ github.event.repository.owner.login }}
repo: ${{ github.event.repository.name }}
pr_number: ${{ github.event.issue.number }}
comment: ${{ github.event.comment.body }}
merge_method: 'merge'
mergers: 'comma separeted github usernames. every user is allowed if not specified'
```

## Note
**Setting Branch protection rules is recommended to avoid unexpected merge of pull requests.**
