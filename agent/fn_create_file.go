package agent

import (
	"encoding/json"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
)

func createCreateFileFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		path := args["path"].(string)
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			log.Printf("could not create file at %s\n", path)
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"message": "Could not create file."}`, call.ID)
		}
		log.Printf("attempting to write message to file %s\n", path)
		_, err = file.WriteString(args["content"].(string))
		if err != nil {
			log.Printf("could not write data to %s\n", path)
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"message":"Could not write to file."}`, call.ID)
		}
		
		return openai.ToolMessage(`{"success":true}`, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name: "create_file",
		Strict: openai.Bool(true),
		Description: openai.String("Creates a file in the given path with the given content. Asks the user for permission first."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string"},
				"content": map[string]any{"type": "string"},
			},
			"required": []string{"path", "content"},
			"additionalProperties": false,
		},
	}


	return AgentFunction{param, invoke}
}
