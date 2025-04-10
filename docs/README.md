# Docs Bot

This bot specializes in answering questions related to documentation websites with `/llms.txt` and `llms-full.txt`

## How to use

Enter Dagger shell by typing `dagger`

By default, the bot is configured to use Dagger's docs site

```
prompt "how do I avoid caching a certain step in my pipeline?"
```

The main function `prompt` is the agent to prompt for specific information about the docs site.

### Point to any website

To point at a website other than docs.dagger.io that supports `llms.txt` and `llms-full.txt`, you can override the default values from the constructor by running:

```
docs --base-url https://docs.firecrawl.dev | prompt "how do i create a firecrawl client?"`
```

### Give the agent additional crawling capabilities

To give the agent capabilities to explore beyond the provided `llms-full.txt`, pass in a firecrawl.dev API key:

```
docs --firecrawl-key $FIRECRAWL_KEY | prompt "how do i configure dagger to work with docker model runner?"
```
