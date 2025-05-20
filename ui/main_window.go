package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
)

type renderableMessage struct {
	sender openai.MessageRole
	content string
}

func CreateMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Galileo")
	w.Resize(fyne.Size{Width: 300, Height: 600})

	messages := []renderableMessage{renderableMessage{openai.MessageRoleAssistant, "Hi there."}}
	
	msgList := widget.NewList(
		func() int {
			return len(messages)
		},
		func() fyne.CanvasObject {
			return CreateSentMessage("", "")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			msg := messages[i]
			var senderText string
			if msg.sender == openai.MessageRoleAssistant {
				senderText = "Galileo"
			} else {
				senderText = "You"
			}
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(senderText)
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(msg.content)
		},
	)
	msgScroll := container.NewScroll(msgList)

	msgInput := widget.NewEntry()
	msgInput.SetPlaceHolder("Enter a message...")
	
	msgSendBtn := widget.NewButtonWithIcon("Send", theme.Icon(theme.IconNameMailSend), func() {
		messages = append(messages, renderableMessage{openai.MessageRoleUser, msgInput.Text})
		msgList.Refresh()
		msgInput.SetText("")
	})
	msgSendBtn.Importance = widget.HighImportance
	
	msgSendPane := container.NewBorder(nil, nil, nil, msgSendBtn, msgInput)
	msgSendPane.Resize(fyne.NewSize(msgSendPane.Size().Width, 50))

	vSplit := container.NewVSplit(msgScroll, msgSendPane)
	vSplit.SetOffset(0.99)

	w.SetContent(vSplit)
	w.ShowAndRun()

	return w
}
