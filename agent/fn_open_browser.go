package agent

import (
	"encoding/json"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
)

// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
func openBrowserToUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func createOpenBrowserFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		url := args["url"].(string)
		openBrowserToUrl(url)
		
		return openai.ToolMessage(`{"success":true}`, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name: "open_browser",
		Strict: openai.Bool(true),
		Description: openai.String("Opens the user's browser to the given URL."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
			"required": []string{"url"},
			"additionalProperties": false,
		},
	}
	
	return AgentFunction{param, invoke}
}
