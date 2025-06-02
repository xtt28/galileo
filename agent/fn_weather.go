package agent

import (
	"log"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
	"github.com/xtt28/galileo/service"

	"encoding/json"
)

func createWeatherFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		apiRes, err := service.GetCurrentWeather()
		if err != nil {
			log.Println("could not fetch weather from service")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"message": "internal error occurred"}`, call.ID)
		}

		resData := apiRes.Properties.Timeseries
		resJson, err := json.Marshal(resData)
		if err != nil {
			log.Println("could not marshal weather to JSON")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"message": "internal error occurred"}`, call.ID)			
		}
		return openai.ToolMessage(string(resJson), call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "get_weather",
		Strict:      openai.Bool(true),
		Description: openai.String("Gets the current weather forecast at the user's location."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{},
			"additionalProperties": false,
		},
	}

	return AgentFunction{param, invoke}
}
