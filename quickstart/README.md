# Dagger Agent Quickstart Experiments

This is an implementation of the [agent quickstart](https://docs.dagger.io/quickstart/agent) with experiments on minimal viable prompts for various popular models.

The prompts can be found under [./prompts](./prompts) titled appropriately for the model name.

The assignment tested with each prompt is:

> write a curl clone

To run this for a specific model, run:

```
go-program "write a curl clone" --model qwen2.5-coder:14b | terminal
```
