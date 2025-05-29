package app

import (
	"fyne.io/fyne/v2/app"
	"github.com/xtt28/galileo/agent"
	"github.com/xtt28/galileo/config"
	"github.com/xtt28/galileo/ui"
)

func Run() {
	conf := config.ReadConfig()

	a := app.New()
	agent.RegisterAllFunctions()

	win := ui.CreateMainWindow(conf.OpenAIKey, a)
	win.Window.ShowAndRun()
}
