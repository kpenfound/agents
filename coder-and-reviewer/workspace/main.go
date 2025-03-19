// A generated module for Workspace functions

package main

import (
	"context"
	"dagger/workspace/internal/dagger"
)

type Workspace struct {
	// +private
	Source *dagger.Directory
}

func New(source *dagger.Directory) *Workspace {
	return &Workspace{
		Source: source,
	}
}

// Read a file in the workspace
func (m *Workspace) Read(
	ctx context.Context,
	// The file path to read
	path string,
) (string, error) {
	return m.Source.File(path).Contents(ctx)
}

// Write a file to a path in the workspace
func (m *Workspace) Write(
	// The file path to write
	path string,
	// The content to write to the file
	content string,
) *Workspace {
	m.Source = m.Source.WithNewFile(path, content)
	return m
}

// Test the workspace
func (m *Workspace) Test(
	ctx context.Context,
	// The test to run
	test string,
) (string, error) {
	tester := dag.Container().From("golang:latest").
		WithWorkdir("/app").
		WithDirectory("/app", m.Source).
		WithExec([]string{"go", "build", "./..."}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})

	exit_code, err := tester.ExitCode(ctx)
	if err != nil {
		return "", err
	}
	if exit_code != 0 {
		return tester.Stderr(ctx)
	}
	return tester.Stdout(ctx)
}

// +internal-use-only DO NOT USE THIS FUNCTION
func (m *Workspace) Work() *dagger.Directory {
	return m.Source
}
