package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
)

func main() {
	var dryRun bool
	var badRunID, badRunURL, badSha, targetSha, targetRunURL, reason, repoPath string
	flag.BoolVar(&dryRun, "dry-run", false, "Print to stdout instead of calling GitHub API")
	flag.StringVar(&badRunID, "bad-run-id", "", "ID of the bad workflow run")
	flag.StringVar(&badRunURL, "bad-run-url", "", "URL of the bad workflow run")
	flag.StringVar(&badSha, "bad-sha", "n/a", "SHA of the bad commit")
	flag.StringVar(&targetSha, "target-sha", "n/a", "SHA of the rollback target commit")
	flag.StringVar(&targetRunURL, "target-run-url", "n/a", "URL of the target workflow run")
	flag.StringVar(&reason, "reason", "", "Reason for rollback")
	flag.StringVar(&repoPath, "repo", os.Getenv("GITHUB_REPOSITORY"), "GitHub repository (owner/repo)")
	flag.Parse()

	if repoPath == "" {
		fmt.Fprintln(os.Stderr, "Error: repository not provided via --repo or GITHUB_REPOSITORY")
		os.Exit(1)
	}

	parts := strings.Split(repoPath, "/")
	owner := parts[0]
	repo := parts[1]

	title := fmt.Sprintf("[auto-rollback] Failed deployment incident run #%s", badRunID)
	body := fmt.Sprintf(`A deployment to GitHub Pages failed post-deploy health verification.

- Bad run: %s
- Bad commit: `+"`%s`"+`
- Rollback target commit: `+"`%s`"+`
- Rollback source run: %s
- Decision: %s`, badRunURL, badSha, targetSha, targetRunURL, reason)

	if dryRun {
		fmt.Printf("DRY RUN: Would create or update issue.\nTitle: %s\nBody:\n%s\n", title, body)
		os.Exit(0)
	}

	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	ctx := context.Background()

	opts := &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list issues: %v\n", err)
		os.Exit(1)
	}

	var existing *github.Issue
	for _, issue := range issues {
		if issue.GetTitle() == title {
			existing = issue
			break
		}
	}

	if existing == nil {
		req := &github.IssueRequest{
			Title:  github.String(title),
			Body:   github.String(body),
			Labels: &[]string{"auto-rollback"},
		}
		issue, _, err := client.Issues.Create(ctx, owner, repo, req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create issue: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created incident issue: %s\n", issue.GetHTMLURL())
		return
	}

	comment := &github.IssueComment{Body: github.String(body)}
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, existing.GetNumber(), comment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create comment: %v\n", err)
		os.Exit(1)
	}

	req := &github.IssueRequest{State: github.String("open")}
	issue, _, err := client.Issues.Edit(ctx, owner, repo, existing.GetNumber(), req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to reopen issue: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated incident issue: %s\n", issue.GetHTMLURL())
}
