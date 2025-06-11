package ai

import (
	"context"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
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

const systemPromptFormat = `You are a helpful assistant and part of the program
Project Galileo. Answer concisely - you are in a conversation with a user. The
current date is %s. Adhere to these guidelines in all your responses:
- For questions related to any current situations, use the web search function.
- If you are opening a file, you may only open files in the user's document folder.
- All files in the user's document folder can be fetched by the respective function.
- Do not open any files or show a message box if the user does not explicitly request
  it.`

func NewConversation(apiKey string) Conversation {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	prompt := fmt.Sprintf(systemPromptFormat, time.Now().String())
	
	conversation := Conversation{
		client,
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.DeveloperMessage(prompt),
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
		dialog.ShowError(err, w)
		log.Println("could not generate chat completion")
		log.Println(err)
		return ""
	}

	choice := completion.Choices[0]
	param := choice.Message.ToParam()
	content := completion.Choices[0].Message.Content
	c.Param.Messages = append(c.Param.Messages, param)

	if len(choice.Message.ToolCalls) > 0 {
		for _, call := range choice.Message.ToolCalls {
			fun, ok := agent.FunctionForName(call.Function.Name)
			if !ok {
				log.Panic("agent attempted to call nonexistent function " + call.Function.Name)
			}
			toolRes := fun.Invoke(w, call)
			content += "\n" + c.SendMessage(w, toolRes)
		}
	}

	return content
}
