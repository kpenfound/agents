// A module for interacting with Documentation websites

package main

import (
	"context"
	"dagger/docs/internal/dagger"
)

type Docs struct {
	BaseURL      string
	Txt          string
	Full         string
	FirecrawlKey *dagger.Secret
}

func New(
	//+default="https://docs.dagger.io"
	baseURL string,
	//+default="llms.txt"
	txt string,
	//+default="llms-full.txt"
	full string,
	//+optional
	firecrawlKey *dagger.Secret,
) Docs {
	return Docs{
		BaseURL:      baseURL,
		Txt:          txt,
		Full:         full,
		FirecrawlKey: firecrawlKey,
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

	if m.FirecrawlKey != nil {
		firecrawl := dag.FirecrawlDag(m.FirecrawlKey)
		env = env.
			WithStringInput("base", m.BaseURL, "the base URL of the documentation website for use with firecrawl").
			WithFirecrawlDagInput("firecrawl", firecrawl, "a tool to use if you cannot find the required information in $llm or $llmsfull. Use $base and paths in $llm to scrape pages. Do not use crawl, only map and scrape.")
	}

	return dag.LLM().
		WithSystemPrompt("You are an agent that can lookup information in documentation and summarize it for a user").
		WithEnv(env).
		WithPrompt(`
			You will be provided a prompt for information about the project
			You have been provided the documentation in the files $llm and $llmsfull
			You can read the contents of the files by selecting them and then using the contents tool. Always do this.
			If the file is too big to understand at once, pass it to the $utils grep tool to search for what youre looking for.
			When searching with the grep tool, use single keywords to search for context.
			Using the files and tools available, answer the prompt as accurately and concicesly as possible, show code examples where applicable
			Keep searching your input files and using your tools until you find the answer
			Your prompt: $prompt`).
		LastReply(ctx)
}
