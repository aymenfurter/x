package main

import (
	"strings"
)

var commandList = []string{
	"nodedisabled",
}

func isCommandInList(cmd string) bool {
	for _, command := range commandList {
		if cmd == command {
			return true
		}
	}
	return false
}

func transformCommand(command string) string {
	segments := strings.Split(command, "&&")
	transformedSegments := make([]string, len(segments))

	for i, segment := range segments {
		trimmedSegment := strings.TrimSpace(segment)

		cmd := strings.Fields(trimmedSegment)[0]
		if isCommandInList(cmd) {
			transformedSegments[i] = "timeout 5 " + trimmedSegment
		} else {
			transformedSegments[i] = trimmedSegment
		}
	}

	return join(transformedSegments)
}

func join(segments []string) string {
	var result strings.Builder

	for i, segment := range segments {
		isLastElement := i == len(segments)-1
		result.WriteString(segment)

		if !isLastElement {

			if strings.HasPrefix(segment, "timeout") {
				result.WriteString(" & ")
			} else {
				result.WriteString(" && ")
			}
		}
	}

	return result.String()
}
