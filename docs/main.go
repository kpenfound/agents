// A module for interacting with Documentation websites

package main

import (
	"context"
)

type Docs struct {
	BaseURL string
	Txt     string
	Full    string
}

func New(
	//+default="https://docs.dagger.io"
	baseURL string,
	//+default="llms.txt"
	txt string,
	//+default="llms-full.txt"
	full string,
) Docs {
	return Docs{
		BaseURL: baseURL,
		Txt:     txt,
		Full:    full,
	}
}

// Returns a container that echoes whatever string argument is provided
func (m *Docs) Prompt(ctx context.Context, prompt string) (string, error) {
	txt := dag.HTTP(m.BaseURL + "/" + m.Txt)
	full := dag.HTTP(m.BaseURL + "/" + m.Full)
	utils := dag.Utils()

	env := dag.Env().
		WithFileInput("llm", txt, "A list of all of the paths in the documentation website and the page descriptions").
		WithFileInput("llmsfull", full, "The entire documentation site expanded as a single markdown file").
		WithStringInput("prompt", prompt, "A prompt for information related to the information contained in the files").
		WithUtilsInput("utils", utils, "Utility tools for searching for content in the docs files")

	return dag.LLM().
		WithEnv(env).
		WithPrompt(`
			You have been provided the documentation for a project in $llm and $llmsfull
			You will be provided a prompt for information about the project
			Using the files and tools available, answer the prompt as accurately and concicesly as possible
			Show code examples where applicable
			Your prompt: $prompt`).
		LastReply(ctx)
}
