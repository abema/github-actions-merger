# github-actions-merger
github-actions-merger is github actions that merges pull request with commit message including pull request labels and release-note block.

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
PullRequest body can include release-note block.

e.g. 
```release-note
Breaking change!
```

3. pull request is merged, and commit message includes labels and release-note block.
```
fix: readme
Labels:
* documentation
* enhancement
release-note:
* Breaking change!
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
enable_auto_merge: true
```

## Options
### Enable Auto Merge
- [About auto merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request#about-auto-merge)
- You can use the auto merge when `enable_auto_merge` is true.
- Default is `false`.
- For more information about enabling auto merge to see the Note: [Enabling auto-merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request#about-auto-merge).


## Note
**Setting Branch protection rules is recommended to avoid unexpected merge of pull requests.**
