package main

import (
	"context"
	"dagger/coding-agent/internal/dagger"
	"fmt"
)

type CodingAgent struct{}

// Write a Go program
func (m *CodingAgent) GoProgram(
	ctx context.Context,
	// The programming assignment, e.g. "write me a curl clone"
	assignment string,
	// Optional prompt file name
	// +optional
	// +default "gpt-4o.md"
	promptFile string,
	// Optional model
	// +optional
	model string,
) (*dagger.Container, error) {
	// Check optional model
	llmopts := dagger.LLMOpts{}
	if model != "" {
		llmopts.Model = model
	}
	// Get the model-specific prompt
	prompt := fmt.Sprintf("prompts/%s", promptFile)

	// Back to the original quickstart
	environment := dag.Env().
		WithStringInput("assignment", assignment, "the assignment to complete").
		WithContainerInput("builder",
			dag.Container().From("golang").WithWorkdir("/app"),
			"a container to use for building Go code").
		WithContainerOutput("completed", "the completed assignment in the Golang container")

	work := dag.LLM(llmopts).
		WithEnv(environment).
		WithPromptFile(dag.CurrentModule().Source().File(prompt))

	return work.
		Env().
		Output("completed").
		AsContainer(), nil
}
