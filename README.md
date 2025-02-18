# Agents

This is a collection of agents that I have created for various projects.

Each of these agents are run with [Dagger](https://dagger.io)

## Try these agents

### Setup

**TEMPORARY EXPERIMENTAL SETUP**
Follow the directions to build the experimental LLM branch of Dagger [here](https://github.com/shykes/melvin)

I'm using ollama over tailscale, so my `.env` looks like:

```bash
OPENAI_BASE_URL=https://my-machine.my-tailnet.ts.net/v1/
OPENAI_MODEL="qwen2.5-coder:32b"
```

If you're using a different service, your `.env` will look different.

### Running the go-programmer agent
From the go-programmer directory, get a Dagger shell. Alternatively, you can load this dagger module remotely with `dagger -m github.com/kpenfound/agents/go-programmer shell`

`_EXPERIMENTAL_DAGGER_RUNNER_HOST=tcp://localhost:1234 ~/bin/dagger-llm shell`

Try the different functions:

This will create a new go project and put you in a terminal in the project's directory:
`assignment "make a hello world app" | terminal`

This will read the [Github issue here](https://github.com/kpenfound/greetings-api/issues/32), read the issue body, complete the work described, and open a PR with the completed work
`solve-issue GITHUB_TOKEN github.com/kpenfound/greetings-api 32`

This will read the [Github pull request here](https://github.com/kpenfound/greetings-api/pull/43), recieve feedback from the input, and add a new commit based on the feedback given.
`feedback GITHUB_TOKEN github.com/kpenfound/greetings-api 43 "please add a test for the new greeting function"`
