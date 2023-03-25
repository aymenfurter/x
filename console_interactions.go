package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func typingEffect(text string, speed time.Duration, color string) {
	colorCode := "\033[38;5;15m"
	if color == "orange" {
		colorCode = "\033[38;5;208m"
	}
	if color == "grey" {
		colorCode = "\033[38;5;242m"
	}

	reset := "\033[0m"

	for _, char := range text {
		fmt.Print(colorCode + string(char) + reset)
		time.Sleep(speed)
	}
}

func askConfirmation() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "y" || choice == "n" || choice == "a" {
			return choice
		}
		fmt.Println("Invalid choice. Please enter Y, N, or A.")
	}
}
