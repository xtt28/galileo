package agent

import (
	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
)

var registeredFunctions = make(map[string]AgentFunction)

type AgentFunction struct {
	APIParams openai.FunctionDefinitionParam
	Invoke func(fyne.Window, openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion
}

func FunctionForName(name string) (fn AgentFunction, ok bool) {
	fn, ok = registeredFunctions[name]
	return
}

func registerFunction(name string, fn AgentFunction) {
	registeredFunctions[name] = fn
}

func GetToolsList() []openai.ChatCompletionToolParam {
	sl := []openai.ChatCompletionToolParam{}
	for _, v := range registeredFunctions {
		param := openai.ChatCompletionToolParam{
			Type: "function",
			Function: v.APIParams,
		}
		sl = append(sl, param)
	}
	return sl
}
	
func RegisterAllFunctions() {
	registerFunction("message_box", createMessageBoxFunc())
	registerFunction("open_browser", createOpenBrowserFunc())
}
