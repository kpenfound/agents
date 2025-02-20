package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"

	github "github.com/google/go-github/v69/github"
)

func main() {
	listener := &Listener{
		WebhookSecretKey: os.Getenv("GH_WEBHOOK_SECRET_KEY"),
	}
	http.Handle("/", http.HandlerFunc(listener.ServeHTTP))
	slog.Info("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

type Listener struct {
	WebhookSecretKey string
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(l.WebhookSecretKey))
	if err != nil {
		slog.Error("Error validating payload", "error", err)
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		slog.Error("Error parsing webhook", "error", err)
	}
	switch event := event.(type) {
	case *github.PushEvent:
		l.ProcessPushEvent(event)
	case *github.IssuesEvent:
		l.ProcessIssuesEvent(event)
	case *github.IssueCommentEvent:
		l.ProcessIssueCommentEvent(event)
	case *github.LabelEvent:
		l.ProcessLabelEvent(event)
	case *github.PullRequestEvent:
		l.ProcessPullRequestEvent(event)
	case *github.PullRequestReviewEvent:
		l.ProcessPullRequestReviewEvent(event)
	}
}

func (l *Listener) ProcessPushEvent(ev *github.PushEvent) {
	slog.Info("Received push event", "repo", ev.Repo.FullName)
	slog.Warn("Unhandled event type", "event", "Push")
}

func (l *Listener) ProcessIssuesEvent(ev *github.IssuesEvent) {
	slog.Info("Received issues event", "repo", ev.Repo.FullName)

	repo := "https://github.com/" + ev.Repo.GetFullName()
	number := ev.Issue.GetNumber()
	assignee := ev.Issue.GetAssignee().GetLogin()

	switch ev.GetAction() {
	case "opened", "edited", "assigned":
		command := fmt.Sprintf("issue --repo %s --issue %d --assignee %s", repo, number, assignee)
		daggerCommand(command)
	default:
		slog.Warn("Unhandled issues event action", "action", ev.GetAction())
	}

}

func (l *Listener) ProcessIssueCommentEvent(ev *github.IssueCommentEvent) {
	slog.Info("Received issue comment event", "repo", ev.Repo.FullName)

	repo := "https://github.com/" + ev.Repo.GetFullName()
	number := ev.Issue.GetNumber()
	comment := ev.Comment.GetBody()
	isPR := ev.Issue.IsPullRequest()

	switch ev.GetAction() {
	case "created":
		command := fmt.Sprintf("issue-comment --repo=%s --issue=%d --comment='%s' --is-pull-request=%t", repo, number, comment, isPR)
		daggerCommand(command)
	default:
		slog.Warn("Unhandled issue comment event action", "action", ev.GetAction())
	}
}

func (l *Listener) ProcessLabelEvent(ev *github.LabelEvent) {
	slog.Info("Received label event", "repo", ev.Repo.FullName)
	slog.Warn("Unhandled event type", "event", "Label")
}

func (l *Listener) ProcessPullRequestEvent(ev *github.PullRequestEvent) {
	slog.Info("Received pull request event", "repo", ev.Repo.FullName)
	slog.Warn("Unhandled event type", "event", "PullRequest")
}

func (l *Listener) ProcessPullRequestReviewEvent(ev *github.PullRequestReviewEvent) {
	slog.Info("Received pull request review event", "repo", ev.Repo.FullName)
	slog.Warn("Unhandled event type", "event", "PullRequestReview")
}

func daggerCommand(command string) {
	slog.Info("Running Dagger command", "command", command)
	cmdSlice := []string{"/usr/local/bin/dagger", "-m", "./dag", "--token", "env:GITHUB_TOKEN", "call", command}
	cmd := exec.Command("sh", "-c", strings.Join(cmdSlice, " "))

	// Capture output and error
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		slog.Error("error running dagger command", "error", fmt.Sprintf("Command failed with error: %v, stdout: %q, stderr: %q", err, out.String(), stderr.String()))
		return
	}

	slog.Info("completed dagger command", "output", out.String())
}
