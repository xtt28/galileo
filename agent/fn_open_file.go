package agent

import (
	"encoding/json"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
)

// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
func openPath(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/C", "start", ""}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func createOpenFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		url := args["path"].(string)
		openPath(url)

		return openai.ToolMessage(`{"success":true}`, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "open",
		Strict:      openai.Bool(true),
		Description: openai.String("Opens either a URL or a file path/app shortcut path in the appropriate application. Before opening any file path or app shortcut, you must first use either get_files or get_programs function to find the exact path of the file that the user needs."),
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
