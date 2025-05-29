package ai

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/xtt28/galileo/agent"
)

const ConversationModel = openai.ChatModelGPT4oMini

type Conversation struct {
	OpenAIClient openai.Client
	Context      context.Context
	Param        openai.ChatCompletionNewParams
}

func NewConversation(apiKey string) Conversation {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	conversation := Conversation{
		client,
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.DeveloperMessage("You are a helpful assistant."),
			},
			Seed:  openai.Int(1),
			Model: ConversationModel,
			Tools: agent.GetToolsList(),
		},
	}
	return conversation
}

func (c *Conversation) SendMessage(w fyne.Window, prompt openai.ChatCompletionMessageParamUnion) string {
	c.Param.Messages = append(c.Param.Messages, prompt)

	completion, err := c.OpenAIClient.Chat.Completions.New(c.Context, c.Param)
	if err != nil {
		// TODO: Improve error handling.
		log.Fatal(err)
	}

	choice := completion.Choices[0]
	param := choice.Message.ToParam()
	content := completion.Choices[0].Message.Content
	c.Param.Messages = append(c.Param.Messages, param)
	
	if len(choice.Message.ToolCalls) > 0 {
		for _, call := range choice.Message.ToolCalls {
			fun, ok := agent.FunctionForName(call.Function.Name)
			if !ok {
				panic("no function with name " + call.Function.Name)
			}
			toolRes := fun.Invoke(w, call)
			content += "\n" + c.SendMessage(w, toolRes)
		}
	}
	

	return content
}
