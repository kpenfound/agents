// A generated module for Coder functions

package main

import (
	"context"
	"dagger/coder/internal/dagger"
)

type Coder struct {
	//+private
	Model string
}

func New(
	model string,
) *Coder {
	return &Coder{
		Model: model,
	}
}

// Solves a coding assignment in a directory
func (m *Coder) Develop(
	ctx context.Context,
	// Assignment or feedback to complete in the directory
	assignment string,
	// source directory containing the project
	source *dagger.Directory,
) *dagger.Directory {
	ws := dag.Workspace(source)

	coder := dag.Llm(dagger.LlmOpts{Model: m.Model}).
		SetString("assignment", assignment).
		SetWorkspace("workspace", ws).
		WithPromptFile(dag.CurrentModule().Source().File("prompt.md"))

	return coder.GetWorkspace("workspace").Work()
}
