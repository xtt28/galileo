package agent

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"

	"encoding/json"
)

func getUserDocumentsPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "Documents")
	default:
		return filepath.Join(os.Getenv("HOME"), "Documents")
	}
}

func getFiles(paths []string) ([]string, error) {
	documentFiles := []string{}
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				documentFiles = append(documentFiles, path)
			}
			return nil
		})
		if err != nil {
			return []string{}, err
		}
	}
	return documentFiles, nil
}

func createShowFilesFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		files, err := getFiles([]string{getUserDocumentsPath()})
		if err != nil {
			log.Println("could not get user documents for agent function")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"An unexpected error occurred."}`, call.ID)
		}

		json, err := json.Marshal(files)
		if err != nil {
			log.Println("could not marshal user documents list to JSON for agent function")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"An unexpected error occurred."}`, call.ID)
		}
		return openai.ToolMessage(string(json), call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "get_files",
		Strict:      openai.Bool(true),
		Description: openai.String("Gives a list of paths to all of the files in the user's documents directory. The document directory is searched recursively."),
		Parameters: openai.FunctionParameters{
			"type":                 "object",
			"properties":           map[string]any{},
			"additionalProperties": false,
		},
	}

	return AgentFunction{param, invoke}
}
