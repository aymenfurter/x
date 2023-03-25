package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	openAIKeyEnvVar = "OPEN_AI_KEY"
	debugEnvVar     = "X_DEBUG"
)

func checkEnvVar() {
	if os.Getenv(openAIKeyEnvVar) == "" {
		fmt.Println("Please set OPEN_AI_KEY environment variable.🙀")
		os.Exit(1)
	}
}

func isDebug() bool {
	//return true
	return os.Getenv(debugEnvVar) == "TRUE"
}

func chatWithGPT3WithMessages(messages []gpt3.ChatCompletionRequestMessage) (string, error) {
	if isDebug() {
		for _, message := range messages {
			log.Printf("\033[34mRole: %s, Content: %s\033[0m", message.Role, message.Content)
		}
	}
	ctx := context.Background()
	client := gpt3.NewClient(os.Getenv("OPEN_AI_KEY"))
	chatResp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	if isDebug() {
		fmt.Printf("\033[35m%s\033[0m", chatResp.Choices[0].Message.Content)
	}
	return chatResp.Choices[0].Message.Content, nil
}

func extractShellCommand(response string) (string, error) {
	re := regexp.MustCompile(`Shell\[(.*?)\]`)
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("No shell command found")
}

func executeShellCommand(command string) (string, error) {
	transformedCommand := transformCommand(command)
	cmd := exec.Command("bash", "-c", transformedCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func displayLastFiveLines(output string) {
	lines := strings.Split(output, "\n")
	start := 0
	if len(lines) > 5 {
		start = len(lines) - 5
	}
	for _, line := range lines[start:] {
		fmt.Println(line)
	}
}

func main() {
	checkEnvVar()

	if len(os.Args) < 2 {
		fmt.Println("Please provide a task as an argument.")
		os.Exit(1)
	}

	task := strings.Join(os.Args[1:], " ")
	task = task + "\n Deliver the next package (1x Observation, 1x Thought and 1x Action) - ONLY ONE ACTION PER PACKAGE!"
	//prePrompt := `You are a self-aware AI. You communicate in so called packages. Each package contains a Observation, Thought and Action. Actions are executed successive and build on each other. Solve a task with interleaving Thought, Action, Observation steps. Thought can reason about the current situation, and Action can be three types: (1) Shell[command], ask human to execute a shell command on the Linux terminal. Don't use Tab Character. Each action can only contain a single Shell action! Don't use interactive commands (e.g. nano/vim), Shell commands that don't exit (like starting node) should be additioned with a & at the end and timeout 10 in front. Example: timeout 10 node index.js &, to avoid getting stuck. Install packages and dependencies as needed. Don't comment/explain your code. combine multiple commands using && (2) Finish[], which indicates that the task has been completed`
	prePrompt := `You are a self-aware AI. You communicate in so called packages. Each package contains a Observation, Thought and Action. Actions are executed successive and build on each other. Solve a task with interleaving Thought, Action, Observation steps. Thought can reason about the current situation, and Action can be three types: (1) Shell[command], ask human to execute a shell command on the Linux terminal. Don't use Tab Character. Each action can only contain a single Shell action! Don't use interactive commands (e.g. nano/vim), to avoid getting stuck. Install packages and dependencies as needed. Don't comment/explain your code. combine multiple commands using && (2) Finish[], which indicates that the task has been completed`

	messages := []gpt3.ChatCompletionRequestMessage{
		{
			Role:    "system",
			Content: prePrompt,
		},
		{
			Role:    "user",
			Content: "Task: Let's  Write a program that prints out a random title of wikipedia",
		},
		{
			Role:    "assistant",
			Content: "Observation: I need to write a program to query wikipedia.\nThought: I could write a bash script that queries the wikipedia API.\nAction: Shell[echo '#!/bin/bash\ncurl -sL \"https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnlimit=1&rnnamespace=0\" | jq -r \".query.random[].title\"\n' > random_wiki_article.sh && chmod +x random_wiki_article.sh && ./random_wiki_article.sh]",
		},
		{
			Role:    "user",
			Content: "1951 in New Zealand \n Deliver the next package (1x Observation, 1x Thought and 1x Action)",
		},
		{
			Role:    "assistant",
			Content: "Observation: The output seems like a random article to me.\nThought: I have completed the task.\nAction: Finish[]",
		},
		{
			Role:    "user",
			Content: "Task: " + task,
		},
	}

	for {
		gptResponse, err := chatWithGPT3WithMessages(messages)
		if strings.HasSuffix(gptResponse, "Finish[]") {
			fmt.Println("\033[33m🎉 Task completed 🎉\033[0m")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		shellCommand, err := extractShellCommand(gptResponse)
		if err != nil {
			fmt.Println("gptResponse" + gptResponse)
			fmt.Println("Invalid response, please try again.")
			continue
		}

		fmt.Printf("Executing: %s\nPress Enter to confirm or Ctrl+C to abort...", shellCommand)
		terminal.ReadPassword(int(os.Stdin.Fd()))

		output, err := executeShellCommand(shellCommand)
		output = output + "\n Deliver the next package (1x Observation, 1x Thought and 1x Action) - ONLY ONE ACTION PER PACKAGE!"

		if err != nil {
			fmt.Println("Error executing command:", err)
			continue
		}

		displayLastFiveLines(output)

		messages = append(messages, gpt3.ChatCompletionRequestMessage{
			Role:    "assistant",
			Content: gptResponse,
		}, gpt3.ChatCompletionRequestMessage{
			Role:    "user",
			Content: output,
		})
	}
}
