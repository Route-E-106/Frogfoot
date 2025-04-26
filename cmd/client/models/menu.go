package model

import (
	"fmt"
	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
    "github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	StateMenu AppState = iota
	StateLogin
	StateRegister
    StateLoggedIn
)

type AppModel struct {
	State     AppState
	MenuIndex int
	Login     Login
	Register  Register
}

func NewApp() AppModel {
	return AppModel{
		State:    StateMenu,
		MenuIndex: 0,
		Login:     NewLogin(),
		Register:  NewRegister(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}

	if _, ok := msg.(utils.BackToMenuMsg); ok {
		m.State = StateMenu
		m.Login = NewLogin() 
		m.Register = NewRegister()
		return m, nil
	}

	switch m.State {
	case StateMenu:
		return m.updateMenu(msg)
	case StateLogin:
		login, cmd := m.Login.Update(msg)
		m.Login = *login
		return m, cmd
	case StateRegister:
		register, cmd := m.Register.Update(msg)
		m.Register = register
		return m, cmd
	}

	return m, nil
}

func (m AppModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up":
			if m.MenuIndex > 0 {
				m.MenuIndex--
			}
		case "down":
			if m.MenuIndex < 1 {
				m.MenuIndex++
			}
		case "enter":
			if m.MenuIndex == 0 {
				m.State = StateLogin
			} else {
				m.State = StateRegister
			}
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m AppModel) View() string {
    
    var boxStyle = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("#7AA2F6")).
        Padding(2, 5)

    var titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#7AA2F6")).
        Padding(0, 1)

    box := "\n" + 
    boxStyle.Render(
        lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("Frogfoot")),
    ) + "\n\n"

	switch m.State {
	case StateMenu:
		cursor := func(i int, text string) string {
            var selectedStyle = lipgloss.NewStyle().
                Foreground(lipgloss.Color("#7AA2F6"))

			if m.MenuIndex == i {
				return selectedStyle.Render("➜ " + text)
			}
			return "  " + text
		}
		return box + fmt.Sprintf(
			"%s\n%s\n\n(Use ↑/↓ and Enter)",
			cursor(0, "Login"), cursor(1, "Register"),
		)
	case StateLogin:
		return box + m.Login.View()
	case StateRegister:
		return box + m.Register.View()
	}
	return ""
}
