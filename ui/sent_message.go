package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
)

func CreateSentMessage(msgData openai.Message) (vbox *fyne.Container) {
	var senderText string
	if msgData.Role == openai.MessageRoleAssistant {
		senderText = "Galileo"
	} else {
		senderText = "You"
	}
	
	senderLabel := widget.NewLabel(senderText)
	senderLabel.TextStyle.Bold = true

	messageContentLabel := widget.NewLabel(msgData.Content[0].Text.Value)
	
	vbox = container.NewVBox(senderLabel, messageContentLabel)
	return
}
