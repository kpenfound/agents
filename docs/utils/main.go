// A utility module for working with documentation site content

package main

import (
	"context"
	"dagger/utils/internal/dagger"
)

type Utils struct{}

// Returns lines that match a pattern in the files of the provided File
func (m *Utils) Grep(
	ctx context.Context,
	// Dagger file to search in
	file *dagger.File,
	// Grep pattern to search in the file
	pattern string,
) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithFile("/mnt/doc.txt", file).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-B", "5", "-A", "5", pattern, "doc.txt"}).
		Stdout(ctx)
}
