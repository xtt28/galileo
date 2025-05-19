package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("Galileo")
	w.Resize(fyne.Size{Width: 300, Height: 600})

	w.SetContent(widget.NewLabel("Welcome to Galileo"))
	w.ShowAndRun()

	return w
}
