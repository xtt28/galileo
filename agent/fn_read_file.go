package agent

import (
	"encoding/json"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/openai/openai-go"
)

func createReadFileFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		path := args["path"].(string)
		called := false
		all := false
		dialog.ShowConfirm("Allow AI to read file?", path, func(allowed bool) {
			called = true
			all = allowed
		}, w)
		for !called {
		}
		if !all {
			log.Printf("user did not allow to read file %s\n", path)
			return openai.ToolMessage(`{"success":false,"description":"User did not allow you to access file."}`, call.ID)
		}
		log.Printf("reading file %s\n", path)
		fileData, err := os.ReadFile(path)
		if err != nil {
			log.Printf("could not read file %s from user's computer.\n", path)
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"Could not access the file."}`, call.ID)
		}

		marshaled, err := json.Marshal(map[string]any{"success": true, "fileData": string(fileData)})
		if err != nil {
			log.Println("could not serialize tool response to JSON")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"Could not serialize to JSON."}`, call.ID)
		}
		return openai.ToolMessage(string(marshaled), call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "read_file",
		Strict:      openai.Bool(true),
		Description: openai.String("Reads the text from the file at the path specified. The user will be asked for permission to read the file before it is read."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
			},
			"required":             []string{"path"},
			"additionalProperties": false,
		},
	}

	return AgentFunction{param, invoke}
}
