Linux tool that generates code using OpenAI's `Codex` model.
You can download the binary for [linux_x64](https://github.com/rm-Umar/codegpt/releases/download/v0.1/codegpt).
# Installation
```go
go build 
```
# Usage
Get your OpenAI API key from [here](https://beta.openai.com/account/api-keys)
and set the env variable.
```bash
export OPENAI_TOKEN=<your_token>
```
```bash
codegpt <instruction>
```
# Example
```bash
$ codegpt 'bash command to remove tmp files'
rm -rf /tmp/*
```
