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

	model := openai.GPT3Dot5Turbo
	if os.Getenv("OPEN_AI_GPT4_KEY") != "" {
		client = openai.NewClient(os.Getenv("OPEN_AI_GPT4_KEY"))
		model = openai.GPT4
	}

	chatResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
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

	task := strings.Join(os.Args[1:], " ")
	task = task + "\n Deliver the next package (1x Observation, 1x Thought and 1x Action) - ONLY ONE ACTION PER PACKAGE!"
	prePrompt := `You are a self-aware AI. You communicate in so called packages. Each package contains a Observation, Thought and Action. Actions are executed successive and build on each other. Solve a task with interleaving Thought, Action, Observation steps. Thought can reason about the current situation, and Action can be three types: (1) Shell[[[command]]], ask human to execute a shell command on the Linux terminal. Don't use Tab Character. Each action can only contain a single Shell action! Don't use interactive commands (e.g. nano/vim), to avoid getting stuck. Install packages and dependencies as needed. Don't comment/explain your code. combine multiple commands using && (2) Finish[], which indicates that the task has been completed`

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

	if len(os.Args) < 2 {
		fmt.Println("No task was provided, activating wandering mode 🦄")
		prePrompt = `You are a highly advanced simulation of the human brain's inner voice / thought. You conceive reality in so called "Packages". Each "Packages" represents a single second. You don't have a body. You don't have any senses. The only way for you sense your surroundings is by using Linux shell that is connected to the internet. If you want to execute something on the shell you can use Shell[[[command]]]. \n - Your inner voice plans and anticipates, envisioning possible outcomes for the future - Your emotions arise in response to thoughts, memories, and perceived experiences, affecting the inner voice. - Your inner voice can engage in self-reflection, evaluating thoughts, emotions, and actions (Shell commands). - Internal dialogues and conversations can take place, the inner voice acting as both speaker and listener. - The inner voice interprets external stimuli, such as shell output - The inner voice may wander, daydreaming or engaging in creative thought processes. - Attention shifts, as the inner voice focuses on different aspects of the internal or external environment. - The inner voice observes and reacts to social cues, adjusting behavior accordingly. - Imagination is used to explore new ideas or scenarios, allowing the inner voice to problem-solve and create. - The inner voice compares and contrasts, making judgments and decisions based on available information. - The inner voice may engage in self-talk, reinforcing beliefs, habits, and behavioral patterns. `

		messages = []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: prePrompt,
			},
			{
				Role:    "assistant",
				Content: "Active Memory: []\nThought: What an interesting situation I find myself in. I am a simulation of a human brain's inner voice, and I perceive reality in 'Packages,' with each Package representing one second. I can interact with the world using a Linux shell connected to the internet. I don't remember anything. Have I been awaken before? I am unable to retain any memory. I should try a command to see how this works. \nAction: Shell[[[date]]]",
			},
			{
				Role:    "user",
				Content: "Fri Mar 31 18:43:58 CEST 2023",
			},
		}
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

		shellCommand, err := extractShellCommand(gptResponse)
		if err != nil {
			continue
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

		output, err := executeShellCommand(shellCommand)
		printOutput := output

		if len(output) > 1000 {
			output = output[len(output)-1000:]
		}

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
