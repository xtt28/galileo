package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"github.com/chromedp/chromedp"
	"github.com/openai/openai-go"
)

func createWebSearchFunc() AgentFunction {
	invoke := func(w fyne.Window, call openai.ChatCompletionMessageToolCall) openai.ChatCompletionMessageParamUnion {
		var args map[string]any
		json.Unmarshal([]byte(call.Function.Arguments), &args)
		query := args["query"].(string)
		log.Printf("attempting to search for query %s\n", query)

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var res string
		err := chromedp.Run(
			ctx,
			chromedp.Navigate(fmt.Sprintf("https://html.duckduckgo.com/html/?q=", url.QueryEscape(query))),
			chromedp.InnerHTML("body", &res, chromedp.NodeVisible),
		)
		if err != nil {
			log.Println("could not search online for query")
			log.Println(err)
			return openai.ToolMessage(`{"success":false,"error":"could not run search"}`, call.ID)
		}
		log.Println("received search query results")

		return openai.ToolMessage(res, call.ID)
	}

	param := openai.FunctionDefinitionParam{
		Name:        "web_search",
		Strict:      openai.Bool(true),
		Description: openai.String("Searches the web for the given query and returns the result page."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
			"additionalProperties": false,
			"required":             []string{"query"},
		},
	}

	return AgentFunction{param, invoke}
}
