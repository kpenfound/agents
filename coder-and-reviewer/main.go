// A generated module for CoderAndReviewer functions

package main

import (
	"context"
	"dagger/coder-and-reviewer/internal/dagger"
)

type CoderAndReviewer struct{}

// Assignment
func (m *CoderAndReviewer) Assignment(
	ctx context.Context,
	// The assignment to complete
	assignment string,
	// The model to use for the assignment
	// +optional
	// +default="gemini-2.0-flash"
	model string,
) *dagger.Directory {
	coder := dag.Coder(model)
	reviewer := dag.Reviewer(model)
	source := dag.Directory()

	work := dag.Llm(dagger.LlmOpts{Model: model}).
		SetString("assignment", assignment).
		SetDirectory("source", source).
		SetCoder("coder", coder).
		SetReviewer("reviewer", reviewer).
		WithPrompt(`
You have an assignment, source directory, coder, and reviewer.
The coder can develop the assignment in the source directory.
The reviewer can review the assignment and give feedback to the coder.
Use the coder to complete work, let the reviewer review the coders work, give the coder feedback from the reviewer, and continue until the reviewer is satisfied
`)
	return work.GetDirectory("source")
}
