package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

var ReadFileTool = openai.FunctionDefinition{
	Name:        "read_file",
	Description: "Read the contents of a file in the working directory. Returns the file's text if it exists, or an error message if it does not.",
	Parameters: map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"filename": map[string]interface{}{
				"type":        "string",
				"description": "The name of the file to read (e.g., 'example.txt').",
			},
		},
		"required": []string{"filename"},
	},
}

func ReadFile(jsonInput string) string {
	// Parse JSON input
	var input struct {
		Filename string `json:"filename"`
	}
	if err := json.Unmarshal([]byte(jsonInput), &input); err != nil {
		return fmt.Sprintf("Error parsing JSON input: %v", err)
	}

	// Validate required field
	if input.Filename == "" {
		return "Error: 'filename' is required"
	}

	// Create full path within workspace
	filePath := filepath.Join("ai_workspace", input.Filename)

	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	return string(content)
}

