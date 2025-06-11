package ui

import (
	"log"

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
	msgBox       *fyne.Container
}

func CreateMainWindow(apiKey string, a fyne.App) MainWindow {
	log.Println("creating main window")
	w := a.NewWindow("Galileo")
	w.Resize(fyne.Size{Width: 300, Height: 600})

	conv := ai.NewConversation(apiKey)
	mw := MainWindow{w, conv, []renderableMessage{}, nil}
	mw.AddWidgets()

	w.ShowAndRun()

	return mw
}

func (mw *MainWindow) AppendMessage(msg renderableMessage) {
	rendered := CreateSentMessage(msg.sender, msg.content)
	mw.Messages = append(mw.Messages, msg)
	mw.msgBox.Add(rendered)
	mw.msgBox.Refresh()
}

func (mw *MainWindow) AddWidgets() {
	log.Println("adding widgets to main window")
	mw.msgBox = container.NewVBox()
	msgScroll := container.NewScroll(mw.msgBox)
	mw.AppendMessage(renderableMessage{openai.MessageRoleAssistant, "Hi there."})

	msgInput := widget.NewEntry()
	msgInput.SetPlaceHolder("Enter a message...")

	msgSendBtn := widget.NewButtonWithIcon("Send", theme.Icon(theme.IconNameMailSend), func() {
		msg := msgInput.Text
		mw.AppendMessage(renderableMessage{openai.MessageRoleUser, msg})
		msgInput.SetText("")

		go func() {
			response := mw.Conversation.SendMessage(mw.Window, openai.UserMessage(msg))
			fyne.Do(func() {
				mw.AppendMessage(renderableMessage{openai.MessageRoleAssistant, response})
				msgScroll.Offset.Y = msgScroll.Content.MinSize().Height - msgScroll.Size().Height
				msgScroll.Base.Refresh()
			})
		}()
	})
	msgSendBtn.Importance = widget.HighImportance

	msgSendPane := container.NewBorder(nil, nil, nil, msgSendBtn, msgInput)
	msgSendPane.Resize(fyne.NewSize(msgSendPane.Size().Width, 50))

	vSplit := container.NewVSplit(msgScroll, msgSendPane)
	vSplit.SetOffset(0.99)
	mw.Window.SetContent(vSplit)
}
