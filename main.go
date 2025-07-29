package main

import (
	"strings"
	"context"
	"fmt"
	"log"
	"bufio"
	"os"
	"flag"

	"github.com/sashabaranov/go-openai"

	"github.com/jakeloud/code/tools"
)

var model = ""
var config = openai.DefaultConfig("ollama")
var client *openai.Client

const system_prompt = `You are a helpful coding assistant`

func init() {
	config.BaseURL = "http://localhost:11434/v1/"
	client = openai.NewClientWithConfig(config)
}

func SendMessage(userMessage string, dialog []openai.ChatCompletionMessage) (string, bool, error) {
	messages := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: system_prompt,
	}}

	for _, message := range dialog {
		messages = append(messages, message)
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Tools: []openai.Tool{
				{
					Type:     openai.ToolTypeFunction,
					Function: &tools.ReadFileTool,
				},
				{
					Type:     openai.ToolTypeFunction,
					Function: &tools.WriteFileTool,
				},
				{
					Type:     openai.ToolTypeFunction,
					Function: &tools.ListFilesTool,
				},
			},
			ToolChoice: "auto",
		},
	)
	if err != nil {
		return "", false, err
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.ToolCalls != nil {
		for _, toolCall := range resp.Choices[0].Message.ToolCalls {
			if toolCall.Type == openai.ToolTypeFunction {
				function_name := toolCall.Function.Name
				arguments := toolCall.Function.Arguments

				fmt.Printf("Using %s with %s\n\n", function_name, arguments)
				
				var output string
				switch (function_name) {
					case "write_file":
						output = tools.WriteFile(arguments)
					case "read_file":
						output = tools.ReadFile(arguments)
					case "list_files":
						output = tools.ListFiles(arguments)

				}
				return output, true, nil
			}
		}
	}

	response := resp.Choices[0].Message.Content
	// this is only relevant to qwen3?
	_, answer, is_thinking := strings.Cut(response, "</think>")
	if is_thinking {
		return strings.TrimSpace(answer), false, nil
	}
	return response, false, nil
}
func List() ([]string, error) {
	result := []string{}
	resp, err := client.ListModels(
		context.Background(),
	)
	if err != nil {
		return result, err
	}
	for _, model := range resp.Models {
		result = append(result, model.ID)
	}
	return result, nil
}

func proompt(input string, history []openai.ChatCompletionMessage) ([]openai.ChatCompletionMessage){
	output, is_tool_call, err := SendMessage(input, history)
	if err != nil {
		log.Fatal(err)
		return history
	}
	if is_tool_call {
		fmt.Printf("Tool call: %s\n", output)
		history = append(history, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: input,
		}, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleTool,
			Content: output,
		})

		return proompt(output, history)
	} else {
		fmt.Printf("%s\n", output)
		history = append(history, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: input,
		}, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant,
			Content: output,
		})
		return history
	}
}

func main() {
	p := flag.String("p", "", "cli prompt")
	flag.Parse()

	history := []openai.ChatCompletionMessage{}
	models, err := List()
	if err != nil {
		log.Fatal(err)
	}
	if len(models) < 1 {
		fmt.Printf("No models ;-(\nExiting...\n")
		return
	}
	model = models[0]

	if p != nil && *p != "" {
		proompt(*p, history)
		return
	}

	fmt.Printf("Models: %v\nUsing: %s\n\n", models, model)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "exit" {
			break
		}
		history = proompt(input, history)
	}
}
