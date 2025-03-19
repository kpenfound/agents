// A generated module for Reviewer functions

package main

import (
	"context"
	"dagger/reviewer/internal/dagger"
)

type Reviewer struct {
	// +private
	Model string
}

func New(
	model string,
) *Reviewer {
	return &Reviewer{
		Model: model,
	}
}

// Returns review of the code changes in the source directory
func (m *Reviewer) Review(
	ctx context.Context,
	// the original assignment being developed
	assignment string,
	// the source directory to review
	source *dagger.Directory,
) (string, error) {
	return dag.Llm(dagger.LlmOpts{Model: m.Model}).
		SetString("assignment", assignment).
		SetDirectory("source", source).
		WithPromptFile(dag.CurrentModule().Source().File("prompt.md")).
		LastReply(ctx)
}
