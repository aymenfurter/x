package main

import (
	"testing"
)

func TestTransformCommandVariousCases(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `echo "This is a test" && node app.js & echo "Done"`,
			expected: `echo "This is a test" && timeout 5 node app.js & echo "Done"`,
		},
		{
			input:    `echo "Line 1 && Line 2" > multiline.txt && node server.js && curl localhost:3000`,
			expected: `echo "Line 1 && Line 2" > multiline.txt && timeout 5 node server.js & curl localhost:3000`,
		},
		{
			input:    `node script.js && echo "Task complete" && echo "Have a nice day"`,
			expected: `timeout 5 node script.js & echo "Task complete" && echo "Have a nice day"`,
		},
		{
			input:    `echo "Starting" && node server.js && sleep 3 && echo "Done"`,
			expected: `echo "Starting" && timeout 5 node server.js & sleep 3 && echo "Done"`,
		},
		{
			input:    `node script1.js && node script2.js & echo "Scripts executed"`,
			expected: `timeout 5 node script1.js & timeout 5 node script2.js & echo "Scripts executed"`,
		},
	}

	for _, tc := range testCases {
		result := transformCommand(tc.input)
		if result != tc.expected {
			t.Errorf("Expected: %s, got: %s", tc.expected, result)
		}
	}
}
