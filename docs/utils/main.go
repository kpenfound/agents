// A utility module for working with documentation site content

package main

import (
	"context"
	"dagger/utils/internal/dagger"
)

type Utils struct{}

// Returns surrounding lines of the file that match a pattern using grep. If a match is not found, returns an error
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
		WithExec([]string{"grep", "-ni", "-C", "10", pattern, "doc.txt"}).
		Stdout(ctx)
}
