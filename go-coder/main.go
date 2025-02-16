// A generated module for GoCoder functions

package main

import (
	"context"
	"dagger/go-coder/internal/dagger"
	"fmt"
)

type GoCoder struct{}

// Ask a go-coder to complete a task and get the Container with the completed task
func (m *GoCoder) Assignment(
	// The task to complete
	task string,
) *dagger.Container {
	// Create a workspace for building Go code
	ws := dag.Workspace(dagger.WorkspaceOpts{
		BaseImage: "golang",
		Context:   dag.Directory(),
		Checker:   "go build ./...",
	})

	// Give the workspace to the LLM
	coder := dag.Llm().
		WithWorkspace(ws).
		WithPromptVar("assignment", task).
		WithPromptFile(dag.CurrentModule().Source().File("prompt.txt")).
		Loop()

	// Return the container
	return coder.Workspace().Container()
}

// Returns a container that echoes whatever string argument is provided
func (m *GoCoder) SolveIssue(
	ctx context.Context,
	// Github authentication token
	githubToken *dagger.Secret,
	// Github repository with an issue to solve
	repo string,
	// Issue number to solve
	issueId int,
) (string, error) {
	// Read assignment from the github issue
	issue := dag.GithubIssue(githubToken).Read(repo, issueId)

	task, err := issue.Body(ctx)
	if err != nil {
		return "", err
	}

	// Complete the assignment
	ws := dag.Workspace(dagger.WorkspaceOpts{
		BaseImage: "golang",
		Context:   dag.Git(repo).Head().Tree(), // Start with the head of the repo
		Checker:   "go build ./...",
	})

	coder := dag.Llm().
		WithWorkspace(ws).
		WithPromptVar("assignment", task).
		WithPromptFile(dag.CurrentModule().Source().File("prompt.txt")).
		Loop()

	completedWork := coder.Workspace().Container().Directory(".")

	// Create a pull request with the completed assignment
	branchName, err := askAnLLM(ctx, task, "Choose a git branch name apprpriate for this assignment. A git branch name should be no more than 20 alphanumeric characters.")
	if err != nil {
		return "", err
	}
	title, err := askAnLLM(ctx, task, "Choose a pull request title that describes the changes made in this assignment.")
	if err != nil {
		return "", err
	}

	featureBranch := dag.FeatureBranch(githubToken, repo, branchName).
		WithChanges(completedWork)

	// Make sure changes have been made to the workspace
	diff, err := featureBranch.Diff(ctx)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("Got empty diff for prompt: %s", task)
	}

	return featureBranch.PullRequest(ctx, title, task)
}

func askAnLLM(ctx context.Context, queryContext, query string) (string, error) {
	return dag.Llm().
		WithPromptVar("query", query).
		WithPromptVar("context", queryContext).
		WithPrompt(`
You will be given a query and a context. Answer the query using the context provided.
Be brief in your responses.
Do not explain the response, simply return the answer.

<query>$query</query>

<context>$context</context>`).
		Loop().
		LastReply(ctx)
}
