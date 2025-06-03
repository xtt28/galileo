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

// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
func getUserApps() ([]string, error) {
	var folders []string
	appFiles := []string{}
	switch runtime.GOOS {
	case "windows":
		folders = []string{
			filepath.Join(os.Getenv("PROGRAMDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
			filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
		}
	}
	for _, folder := range folders {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				log.Println("getting installed app: " + path)
				appFiles = append(appFiles, path)
			}
			return nil
		})
		if err != nil {
			log.Println(err)
			return appFiles, err
		}
	}
	return appFiles, nil
}

func createGetProgramsFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		apps, err := getUserApps()
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
