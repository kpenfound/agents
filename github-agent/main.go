// A generated module for GithubAgent functions

package main

import (
	"context"
	"dagger/github-agent/internal/dagger"
	"strings"
)

const (
	DAGGER_LLM_VERSION = "0.17.1"
	BOT_NAME           = "kpenfound"
)

type GithubAgent struct {
	GithubToken *dagger.Secret
	Model       string
}

func New(
	// GitHub Token
	token *dagger.Secret,
	// Model to use for solving issues
	// +optional
	// +default="qwen2.5-coder:32b"
	model string,
) *GithubAgent {
	return &GithubAgent{
		GithubToken: token,
		Model:       model,
	}
}

// Start the GithubAgent Webhook Listener
func (m *GithubAgent) Listen(webhookSecret, daggerCloudToken, openAiBaseUrl *dagger.Secret) *dagger.Service {
	me := dag.CurrentModule().Source()
	return dag.Container().
		From("golang:alpine").
		WithEnvVariable("DAGGER_VERSION", DAGGER_LLM_VERSION).
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"sh", "-c", "curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=/usr/local/bin sh"}).
		WithSecretVariable("GH_WEBHOOK_SECRET_KEY", webhookSecret).
		WithSecretVariable("GITHUB_TOKEN", m.GithubToken).
		WithSecretVariable("DAGGER_CLOUD_TOKEN", daggerCloudToken).
		WithSecretVariable("OPENAI_BASE_URL", openAiBaseUrl).
		WithWorkdir("/src").
		WithDirectory("/src", me.Directory("webhook-server")).
		WithDirectory("/src/dag", me.WithoutDirectory("webhook-server")).
		WithExposedPort(8080).
		AsService(dagger.ContainerAsServiceOpts{
			Args:                          []string{"go", "run", "."},
			ExperimentalPrivilegedNesting: true,
		})
}

func (m *GithubAgent) IssueComment(ctx context.Context, repo string, issue int, comment string, isPullRequest bool) (string, error) {
	// Comment was addressed to the Bot
	if isPullRequest && strings.HasPrefix(comment, BOT_NAME) {
		return dag.GoCoder().PrFeedback(ctx, m.GithubToken, repo, issue, comment, dagger.GoCoderPrFeedbackOpts{Model: m.Model})

	}
	return "NOOP", nil
}

func (m *GithubAgent) Issue(ctx context.Context, repo string, issue int, assignee string) (string, error) {
	if assignee == BOT_NAME {
		return dag.GoCoder().SolveIssue(ctx, m.GithubToken, repo, issue, dagger.GoCoderSolveIssueOpts{Model: m.Model})
	}
	return "NOOP", nil
}
