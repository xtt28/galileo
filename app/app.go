package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xtt28/galileo/config"
	"github.com/xtt28/galileo/ui"

	"log"
)

func Run() {
	config.ReadConfig()
	prog := tea.NewProgram(ui.InitialModel())
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
