package agent

import (
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/openai/openai-go"
)

func createMessageBoxFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)

		dialog.ShowInformation("", args["message"].(string), w)
		return openai.ToolMessage(`{"success":true}`, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "message_box",
		Strict:      openai.Bool(true),
		Description: openai.String("Shows a GUI message box to the user with the given text. Do not use unless explicitly told to do so."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{"type": "string"},
			},
			"required":             []string{"message"},
			"additionalProperties": false,
		},
	}

	return AgentFunction{param, invoke}
}
