package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func summerize(inputText string) (string, error) {
	prePrompt := `Describe this shell output in 50 characters or less, the summery should start with "The command returned.. ":`

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: prePrompt,
		},
		{
			Role:    "user",
			Content: "Describe this shell output in 50 characters or less, the summery should start with \"The command returned.. \"" + inputText,
		},
	}

	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPEN_AI_KEY"))
	chatResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		fmt.Printf("\033[35m%s\033[0m", chatResp.Choices[0].Message.Content)
	}

	return chatResp.Choices[0].Message.Content, nil
}
