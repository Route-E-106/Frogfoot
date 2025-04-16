// based on https://github.com/charmbracelet/bubbletea-app-template
package main

import (
	"fmt"
	"os"
	"github.com/Route-E-106/Frogfoot/cmd/client/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
