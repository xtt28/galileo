package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
)

func CreateSentMessage(role openai.MessageRole, content string) (vbox *fyne.Container) {
	var senderText string
	if role == openai.MessageRoleAssistant {
		senderText = "Galileo"
	} else {
		senderText = "You"
	}
	
	senderLabel := widget.NewLabel(senderText)
	senderLabel.TextStyle.Bold = true

	messageContentLabel := widget.NewLabel(content)
	messageContentLabel.Wrapping = fyne.TextWrapWord
	
	vbox = container.NewVBox(senderLabel, messageContentLabel)
	return
}
