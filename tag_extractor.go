package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func extractTags(inputText string) (string, error) {
	prePrompt := `Extract a comma separated list of relevant tags (dont use #) for the following text:`

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: prePrompt,
		},
		{
			Role:    "user",
			Content: "Extract a comma separated list of relevant tags (dont use #) for the following text:" + inputText,
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
