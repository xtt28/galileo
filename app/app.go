package app

import (
	"fyne.io/fyne/v2/app"
	"github.com/xtt28/galileo/config"
	"github.com/xtt28/galileo/ui"

	"log"
)

func Run() {
	config.ReadConfig()	

	a := app.New()
	ui.CreateMainWindow(a).ShowAndRun()
}
