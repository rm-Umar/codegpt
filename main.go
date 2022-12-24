package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type TextCompletion struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int            `json:"created"`
	Model   string         `json:"model"`
	Choices []Choice       `json:"choices"`
	Usage   map[string]int `json:"usage"`
}

type Choice struct {
	Text          string      `json:"text"`
	Index         int         `json:"index"`
	Logprobs      interface{} `json:"logprobs"`
	Finish_reason string      `json:"finish_reason"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		printHelp()
		return
	}
	prompt := args[1]
	apiToken := os.Getenv("OPENAI_TOKEN")
	if apiToken == "" {
		fmt.Println("Please set the env. var.\nexport OPENAI_TOKEN=<Your Token>")
		return
	}

	data := []byte(`{
		"model": "text-davinci-003", 
		"prompt": "` + prompt + `" , 
		"temperature": 0, 
		"max_tokens": 2048,
		"top_p": 1
		}`)
	r := bytes.NewReader(data)

	// Set up the POST request
	url := "https://api.openai.com/v1/completions"
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request Failed Try again")
		return
	}
	var textResponse TextCompletion
	err = json.Unmarshal(body, &textResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Use the Go lexer to tokenize the code
	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, textResponse.Choices[0].Text)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Use the solarized style for syntax highlighting
	style := styles.Get("solarized-dark")

	// Use the terminal formatter to output the highlighted code
	formatter := formatters.Get("terminal")
	var writer bytes.Buffer
	if err := formatter.Format(&writer, style, iterator); err != nil {
		fmt.Println(err)
	}
	fmt.Println(writer.String())
}

func printHelp() {
	fmt.Printf(`
	Usage: 
		codegpt <instruction>
	Example:
		help 'bash command to remove tmp files'
		rm -rf /tmp/*
	`)
}
