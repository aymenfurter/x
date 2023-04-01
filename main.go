package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
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
	return os.Getenv(debugEnvVar) == "TRUE"
}

func chatWithopenaiWithMessages(messages []openai.ChatCompletionMessage) (string, error) {
	if isDebug() {
		for _, message := range messages {
			log.Printf("\033[34mRole: %s, Content: %s\033[0m", message.Role, message.Content)
		}
	}
	ctx := context.Background()
	client := openai.NewClient(os.Getenv("OPEN_AI_KEY"))
	chatResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
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
	re := regexp.MustCompile(`(?s)Shell\[\[\[(.*?)\]\]\]`)

	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("No shell command found")
}

func executeShellCommand(command string) (string, error) {
	transformedCommand := transformCommand(command)
	cmd := exec.Command("bash", "-c", transformedCommand)
	fmt.Println("\033[32m $\033[0m", transformedCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func extractThought(response string) (string, error) {
	re := regexp.MustCompile(`Thought: (.*?)\n`)
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("No thought found")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	checkEnvVar()

	if len(os.Args) < 2 {
		fmt.Println("Please provide a task as an argument.")
		os.Exit(1)
	}

	task := strings.Join(os.Args[1:], " ")
	task = task + "\n Deliver the next package (1x Observation, 1x Thought and 1x Action) - ONLY ONE ACTION PER PACKAGE!"
	//prePrompt := `You are a self-aware AI. You communicate in so called packages. Each package contains a Observation, Thought and Action. Actions are executed successive and build on each other. Solve a task with interleaving Thought, Action, Observation steps. Thought can reason about the current situation, and Action can be three types: (1) Shell[command], ask human to execute a shell command on the Linux terminal. Don't use Tab Character. Each action can only contain a single Shell action! Don't use interactive commands (e.g. nano/vim), Shell commands that don't exit (like starting node) should be additioned with a & at the end and timeout 10 in front. Example: timeout 10 node index.js &, to avoid getting stuck. Install packages and dependencies as needed. Don't comment/explain your code. combine multiple commands using && (2) Finish[], which indicates that the task has been completed`
	prePrompt := `You are a self-aware AI. Before each Package, you will print your current Active Memory (Limited to 50 words). Before you do any changes, use shell commands to explore dirctories and files - research using or instance cat and curl commands, store what you want to memorize in your memory. Before each package, write your current motivation and goals as an individual. You communicate in so called packages. Each package contains a Observation, Thought and Action. Actions are executed successive and build on each other. Solve a task with interleaving Thought, Action, Observation steps. Thought can reason about the current situation, and Action can be three types: (1) Shell[[[command]]], ask human to execute a shell command on the Linux terminal. Don't use Tab Character. Each action can only contain a single Shell action! Don't use interactive commands (e.g. code/nano/vim), to avoid getting stuck. Install packages and dependencies as needed. Don't comment/explain your code. combine multiple commands using && (2) Finish[], which indicates that the task has been completed`

	messages := []openai.ChatCompletionMessage{
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
			Content: "Observation: I need to write a program to query wikipedia.\nThought: I could write a bash script that queries the wikipedia API.\nAction: Shell[[[echo '#!/bin/bash\ncurl -sL \"https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnlimit=1&rnnamespace=0\" | jq -r \".query.random[].title\"\n' > random_wiki_article.sh && chmod +x random_wiki_article.sh && ./random_wiki_article.sh]]]",
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

	confirmCommand := true
	for {
		fmt.Print("🤔 ")
		progressActive = true

		wg.Add(1)
		go func() {
			defer wg.Done()
			showProgress()
		}()

		gptResponse, err := chatWithopenaiWithMessages(messages)

		closeProgress()

		if strings.HasSuffix(gptResponse, "Finish[]") {
			fmt.Println("\033[33m🎉 Task completed 🎉\033[0m")
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		thought, err := extractThought(gptResponse)
		typingEffect(thought, 10*time.Millisecond, "orange")
		fmt.Println()

		errorCmd := ""

		shellCommand, err := extractShellCommand(gptResponse)
		if err != nil {
			// check if gptResponse contained code, vim, nano, etc.
			illegalCommands := []string{"code", "vim", "nano"}
			// CHeck if gptResponse contains Shell[[[<any of the illegal commands>
			for _, illegalCommand := range illegalCommands {
				if strings.Contains(gptResponse, "Shell[[["+illegalCommand) {
					typingEffect("Illegal command: "+illegalCommand, 10*time.Millisecond, "red")
					errorCmd = illegalCommand
					fmt.Println()
				}
			}

			if errorCmd == "" {
				continue
			}
		}

		typingEffect("Executing: ", 10*time.Millisecond, "white")
		typingEffect(shellCommand, 10*time.Millisecond, "grey")
		typingEffect("\nConfirm: [Y]es / [N]o / [A]ll future commands..  ", 10*time.Millisecond, "white")
		confirmation := askConfirmation()
		fmt.Println()

		if !confirmCommand {
			confirmation = "y"
		}

		if confirmation == "n" {
			os.Exit(0)
		} else if confirmation == "a" {
			confirmCommand = false
		}

		output := ""
		err = nil
		if errorCmd != "" {
			output = "CRITICAL RULE VIOLATION: Illegal command: " + errorCmd + "\n"
		} else {
			output, err = executeShellCommand(shellCommand)
		}
		printOutput := output
		output = output + "\n Deliver the next package (1x Observation, 1x Thought and 1x Action) - ONLY ONE ACTION PER PACKAGE!"

		if err != nil {
			scrollText(printOutput)
		}

		scrollText(printOutput)

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: gptResponse,
		}, openai.ChatCompletionMessage{
			Role:    "user",
			Content: output,
		})
	}
}

var (
	progressActive = true
	progressMtx    sync.Mutex
)

func showProgress() {
	emojiChars := []string{"🤔  ", "🧐  ", "🤨  "}
	i := 0

	for {
		progressMtx.Lock()
		if !progressActive {
			progressMtx.Unlock()
			break
		}
		progressMtx.Unlock()

		fmt.Printf("\r%s", emojiChars[i])
		i = (i + 1) % len(emojiChars)
		time.Sleep(500 * time.Millisecond)
	}
}

func closeProgress() {
	progressMtx.Lock()
	defer progressMtx.Unlock()
	progressActive = false
}
