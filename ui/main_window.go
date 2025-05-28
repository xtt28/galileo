package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"
	"github.com/xtt28/galileo/ai"
)

type renderableMessage struct {
	sender  openai.MessageRole
	content string
}

type MainWindow struct {
	Window       fyne.Window
	Conversation ai.Conversation
	Messages     []renderableMessage
}

func CreateMainWindow(apiKey string, a fyne.App) MainWindow {
	w := a.NewWindow("Galileo")
	w.Resize(fyne.Size{Width: 300, Height: 600})

	conv := ai.NewConversation(apiKey)
	mw := MainWindow{w, conv, []renderableMessage{renderableMessage{openai.MessageRoleAssistant, "hi there"}}}
	mw.AddWidgets()

	w.ShowAndRun()

	return mw
}

func (mw *MainWindow) AddWidgets() {
	msgList := widget.NewList(
		func() int {
			return len(mw.Messages)
		},
		func() fyne.CanvasObject {
			return CreateSentMessage("", "")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			msg := mw.Messages[i]
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
		msg := msgInput.Text
		mw.Messages = append(mw.Messages, renderableMessage{openai.MessageRoleUser, msg})
		msgList.Refresh()
		msgInput.SetText("")

		fyne.Do(func() {
			response := mw.Conversation.SendMessage(openai.UserMessage(msg))
			mw.Messages = append(mw.Messages, renderableMessage{openai.MessageRoleAssistant, response})
			msgList.Refresh()
		})

	})
	msgSendBtn.Importance = widget.HighImportance

	msgSendPane := container.NewBorder(nil, nil, nil, msgSendBtn, msgInput)
	msgSendPane.Resize(fyne.NewSize(msgSendPane.Size().Width, 50))

	vSplit := container.NewVSplit(msgScroll, msgSendPane)
	vSplit.SetOffset(0.99)
	mw.Window.SetContent(vSplit)
}
