package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
)

func writeOutput(outputs map[string]string) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		for k, v := range outputs {
			fmt.Printf("%s=%s\n", k, v)
		}
	} else {
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to GITHUB_OUTPUT: %v\n", err)
		} else {
			defer f.Close()
			for k, v := range outputs {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}
		}
	}

	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile != "" {
		f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			badSha := outputs["bad_sha"]
			if badSha == "" {
				badSha = "n/a"
			}
			badUrl := outputs["bad_run_url"]
			if badUrl == "" {
				badUrl = "n/a"
			}
			targetSha := outputs["target_sha"]
			if targetSha == "" {
				targetSha = "n/a"
			}
			targetUrl := outputs["target_run_url"]
			if targetUrl == "" {
				targetUrl = "manual"
			}

			summary := fmt.Sprintf(`## Auto Rollback Decision
- Should rollback: `+"`%s`"+`
- Incident detected: `+"`%s`"+`
- Reason: %s
- Bad commit: `+"`%s`"+`
- Bad run: %s
- Target commit: `+"`%s`"+`
- Target run: %s
`, outputs["should_rollback"], outputs["incident_detected"], outputs["reason"], badSha, badUrl, targetSha, targetUrl)
			fmt.Fprintln(f, summary)
		}
	}
}

func main() {
	outputs := map[string]string{
		"should_rollback":   "false",
		"incident_detected": "false",
		"reason":            "",
		"bad_sha":           "",
		"bad_run_id":        "",
		"bad_run_url":       "",
		"target_sha":        "",
		"target_run_id":     "",
		"target_run_url":    "",
	}

	exitWith := func(reason string) {
		outputs["reason"] = reason
		writeOutput(outputs)
		os.Exit(0)
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	repoPath := os.Getenv("GITHUB_REPOSITORY")
	if repoPath == "" {
		exitWith("GITHUB_REPOSITORY not set")
		return
	}
	parts := strings.Split(repoPath, "/")
	owner := parts[0]
	repo := parts[1]

	if eventName == "workflow_dispatch" {
		targetSha := os.Getenv("INPUT_ROLLBACK_SHA")
		badSha := os.Getenv("INPUT_BAD_SHA")
		outputs["should_rollback"] = "true"
		outputs["target_sha"] = targetSha
		outputs["bad_sha"] = badSha
		exitWith("Manual rollback requested via workflow_dispatch.")
		return
	}

	autoRollbackEnabled := strings.ToLower(os.Getenv("AUTO_ROLLBACK_ENABLED"))
	if autoRollbackEnabled == "false" {
		exitWith("Auto rollback disabled by repository variable AUTO_ROLLBACK_ENABLED=false.")
		return
	}

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		exitWith(fmt.Sprintf("Failed to read GITHUB_EVENT_PATH: %v", err))
		return
	}

	var payload struct {
		WorkflowRun *github.WorkflowRun `json:"workflow_run"`
	}
	if err := json.Unmarshal(eventData, &payload); err != nil {
		exitWith(fmt.Sprintf("Failed to parse event JSON: %v", err))
		return
	}
	run := payload.WorkflowRun
	if run == nil {
		exitWith("No workflow_run found in event payload.")
		return
	}

	outputs["bad_sha"] = run.GetHeadSHA()
	outputs["bad_run_id"] = fmt.Sprintf("%d", run.GetID())
	outputs["bad_run_url"] = run.GetHTMLURL()

	if run.GetConclusion() != "failure" {
		exitWith(fmt.Sprintf("Triggering workflow conclusion '%s' is not eligible for rollback.", run.GetConclusion()))
		return
	}
	if run.GetRunAttempt() > 1 {
		exitWith(fmt.Sprintf("Skipping repeated attempt for incident run %d (run_attempt=%d).", run.GetID(), run.GetRunAttempt()))
		return
	}

	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	ctx := context.Background()

	opts := &github.ListWorkflowJobsOptions{ListOptions: github.ListOptions{PerPage: 100}}
	jobs, _, err := client.Actions.ListWorkflowJobs(ctx, owner, repo, run.GetID(), opts)
	if err != nil {
		exitWith(fmt.Sprintf("Failed to fetch jobs: %v", err))
		return
	}

	var deployJob *github.WorkflowJob
	for _, j := range jobs.Jobs {
		if j.GetName() == "merge_and_deploy" {
			deployJob = j
			break
		}
	}

	if deployJob == nil {
		exitWith("Deployment job 'merge_and_deploy' not found; no auto rollback.")
		return
	}

	postDeployFailed := false
	for _, step := range deployJob.Steps {
		if step.GetName() == "Post-deploy Health Verification" && step.GetConclusion() == "failure" {
			postDeployFailed = true
			break
		}
	}

	if !postDeployFailed {
		exitWith("Deployment did not fail due to post-deploy health verification; no auto rollback.")
		return
	}

	outputs["incident_detected"] = "true"

	runsOpts := &github.ListWorkflowRunsOptions{
		Branch:      "main",
		Status:      "completed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	successfulRuns, _, err := client.Actions.ListWorkflowRunsByID(ctx, owner, repo, run.GetWorkflowID(), runsOpts)
	if err != nil {
		exitWith(fmt.Sprintf("Failed to fetch workflow runs: %v", err))
		return
	}

	badRunCreatedAt := run.GetCreatedAt().Time
	var priorGoodRun *github.WorkflowRun

	for _, candidate := range successfulRuns.WorkflowRuns {
		if candidate.GetConclusion() == "success" &&
			candidate.GetHeadSHA() != "" &&
			candidate.GetHeadSHA() != run.GetHeadSHA() &&
			candidate.GetCreatedAt().Time.Before(badRunCreatedAt) {
			priorGoodRun = candidate
			break
		}
	}

	if priorGoodRun == nil {
		exitWith("No prior known-good successful deployment run found.")
		return
	}
	if priorGoodRun.GetHeadSHA() == run.GetHeadSHA() {
		exitWith("Resolved rollback target equals failing commit; skipping to prevent rollback loop.")
		return
	}

	outputs["should_rollback"] = "true"
	outputs["target_sha"] = priorGoodRun.GetHeadSHA()
	outputs["target_run_id"] = fmt.Sprintf("%d", priorGoodRun.GetID())
	outputs["target_run_url"] = priorGoodRun.GetHTMLURL()
	exitWith(fmt.Sprintf("Auto rollback target resolved from run %d.", priorGoodRun.GetID()))
}
