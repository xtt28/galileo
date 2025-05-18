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
	Messages []openai.ChatCompletionMessageParamUnion
}

func NewConversation(apiKey string) Conversation {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	conversation := Conversation{
		client,
		context.Background(),
		[]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful chatbot."),
		},
	}
	return conversation
}

func (c *Conversation) SendMessage(prompt string) {
	message := openai.UserMessage(prompt)
	c.Messages = append(c.Messages, message)
	
	params := openai.ChatCompletionNewParams{
		Messages: c.Messages,
		Model: ConversationModel,
	}
	
	completion, err := c.OpenAIClient.Chat.Completions.New(c.Context, params)
	if err != nil {
		log.Fatal(err)
	}

	response := completion.Choices[0].Message.ToParam()
	c.Messages = append(c.Messages, response)
}
