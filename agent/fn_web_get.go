package agent

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"github.com/chromedp/chromedp"
	"github.com/openai/openai-go"
)

func createWebGetFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		urlVisit := args["url"].(string)
		log.Printf("attempting to access url %s\n", urlVisit)

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		ctx, cancel2 := chromedp.NewContext(ctxWithTimeout)
		defer cancel2()

		var res string
		err := chromedp.Run(ctx,
			chromedp.Navigate(urlVisit),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.Text("body", &res, chromedp.ByQuery),
		)

		if err != nil {
			log.Printf("could not search url %s online\n", urlVisit)
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"error":"could not access url"}`, call.ID)
		}
		log.Println("received url results")
		
		return openai.ToolMessage(res, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name: "web_get",
		Strict: openai.Bool(true),
		Description: openai.String("Gets the text from a webpage at the given address."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{"type": "string"},
			},
			"additionalProperties": false,
			"required": []string{"url"},
		},
	}

	return AgentFunction{param, invoke}
}
