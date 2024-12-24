# github-actions-merger

github-actions-merger is a custom GitHub Action that merges pull request with metadata (commit message, pull request labels, release-note block, and git trailers).

## Usage

Write your workflow file.

```yaml
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

Post a comment with ```/merge``` on a GitHub pull request.

A pull-request body can include release-note block.

e.g.

```release-note
Breaking change!
```

The pull request will be merged, and commit message includes labels and release-note block as following.

~~~md
fix: readme
---
Labels:
* documentation
* enhancement
```release-note
Breaking change!
```
~~~


## Parameters

You need to set parameters in workflow.

```yaml
github_token: ${{ secrets.GITHUB_TOKEN }}
owner: ${{ github.event.repository.owner.login }}
repo: ${{ github.event.repository.name }}
pr_number: ${{ github.event.issue.number }}
comment: ${{ github.event.comment.body }}
merge_method: 'merge'
mergers: 'comma separeted github usernames. every user is allowed if not specified'
enable_auto_merge: true
git_trailers: 'Co-authored-by=abema,Co-authored-by=actions'
```


## Options

### Enable Auto Merge

- [About auto merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request#about-auto-merge)
- You can use the auto merge when `enable_auto_merge` is true.
- Default is `false`.
- For more information about enabling auto merge to see the Note: [Enabling auto-merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request#about-auto-merge).


## Note

**Setting the branch protection rules is recommended to avoid unexpected merging of pull requests.**
