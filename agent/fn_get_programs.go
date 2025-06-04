package agent

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
)

func getProgramShortcutPaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			filepath.Join(os.Getenv("PROGRAMDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
			filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
		}
	}
	return []string{}
}

func createGetProgramsFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		apps, err := getFiles(getProgramShortcutPaths())
		if err != nil {
			log.Println("could not get user apps for agent function")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"An unexpected error occurred."}`, call.ID)
		}

		json, err := json.Marshal(apps)
		if err != nil {
			log.Println("could not marshal user app list to JSON for agent function")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"description":"An unexpected error occurred."}`, call.ID)
		}
		return openai.ToolMessage(string(json), call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "get_apps",
		Strict:      openai.Bool(true),
		Description: openai.String("Gets a list of apps installed on the user's computer as a list of paths to shortcut files on the computer."),
		Parameters: openai.FunctionParameters{
			"type":                 "object",
			"properties":           map[string]any{},
			"additionalProperties": false,
		},
	}

	return AgentFunction{param, invoke}
}
