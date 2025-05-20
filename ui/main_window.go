package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/openai/openai-go"

	"log"
)

func CreateMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Galileo")
	w.Resize(fyne.Size{Width: 300, Height: 600})

	msgList := container.NewScroll(CreateSentMessage(openai.UserMessage("hi there")))

	msgInput := widget.NewEntry()
	msgInput.SetPlaceHolder("Enter a message...")
	
	msgSendBtn := widget.NewButtonWithIcon("Send", theme.Icon(theme.IconNameMailSend), func() {
		log.Println("Message: " + msgInput.Text)
		msgInput.SetText("")
	})
	msgSendBtn.Importance = widget.HighImportance
	
	msgSendPane := container.NewBorder(nil, nil, nil, msgSendBtn, msgInput)
	msgSendPane.Resize(fyne.NewSize(msgSendPane.Size().Width, 50))

	vSplit := container.NewVSplit(msgList, msgSendPane)
	vSplit.SetOffset(0.99)

	w.SetContent(vSplit)
	w.ShowAndRun()

	return w
}
