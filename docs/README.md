# Docs Bot

This bot specializes in answering questions related to documentation websites with `/llms.txt` and `llms-full.txt`

## How to use

Enter Dagger shell by typing `dagger`

By default, the bot is configured to use Dagger's docs site

```
prompt "how do I avoid caching a certain step in my pipeline?"
```

The main function `prompt` is the agent to prompt for specific information about the docs site.

To point at a documentation website other than docs.dagger.io, you can override the default values from the constructor by running:

```
docs --base-url https://docs.firecrawl.dev | prompt "how do i create a firecrawl client?"`
```
