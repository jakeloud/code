package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/sashabaranov/go-openai"
	"os"
	"path/filepath"
)

var WriteFileTool = openai.FunctionDefinition{
	Name:        "write_file",
	Description: "Write or overwrite content to a file in the working directory. Creates the file if it doesn't exist, and returns a confirmation message.",
	Parameters: map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"filename": map[string]interface{}{
				"type":        "string",
				"description": "The name of the file to write (e.g., 'example.txt').",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The text content to write into the file.",
			},
		},
		"required": []string{"filename", "content"},
	},
}

func WriteFile(jsonInput string) string {
	// Parse JSON input
	var input struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}
	if err := json.Unmarshal([]byte(jsonInput), &input); err != nil {
		return fmt.Sprintf("Error parsing JSON input: %v", err)
	}

	// Validate required fields
	if input.Filename == "" {
		return "Error: 'filename' is required"
	}
	if input.Content == "" {
		return "Error: 'content' is required"
	}

	// Create ai_workspace directory if it doesn't exist
	workspace := "ai_workspace"
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	// Create absolute path within workspace
	filePath := filepath.Join(workspace, input.Filename)

	// Write file with content
	if err := ioutil.WriteFile(filePath, []byte(input.Content), 0644); err != nil {
		return fmt.Sprintf("Error writing to file: %v", err)
	}

	return fmt.Sprintf("Successfully wrote %d bytes to %s", len(input.Content), filePath)
}
