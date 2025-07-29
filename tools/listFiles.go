package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var ListFilesTool = openai.FunctionDefinition{
	Name:        "list_files",
	Description: "List the names of all files in the working directory. Returns a newline-separated string of filenames, or 'No files' if the directory is empty.",
	Parameters: map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{},
		"required": []string{},
	},
}

func ListFiles(jsonInput string) string {
	// Define the workspace directory
	workspace := "ai_workspace"

	// Create the workspace directory if it doesn't exist
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	// Read the directory contents
	files, err := ioutil.ReadDir(workspace)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	// Filter out directories and collect file names
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	// Handle empty directory case
	if len(fileNames) == 0 {
		return "No files in ai_workspace"
	}

	// Return newline-separated list of files
	return strings.Join(fileNames, "\n")
}

