package ai

import (
	"context"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const ConversationModel = openai.ChatModelGPT4oMini

type Conversation struct {
	OpenAIClient openai.Client
	Context context.Context
	Param openai.ChatCompletionNewParams
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
			Seed: openai.Int(1),
			Model: ConversationModel,
		},
	}
	return conversation
}

func (c *Conversation) SendMessage(prompt openai.ChatCompletionMessageParamUnion) {
	c.Param.Messages = append(c.Param.Messages, prompt)

	completion, err := c.OpenAIClient.Chat.Completions.New(c.Context, c.Param)
	if err != nil {
		// TODO: Improve error handling.
		log.Fatal(err)
	}

	// TODO: Handle tool calls.

	c.Param.Messages = append(c.Param.Messages, completion.Choices[0].Message.ToParam())
}
