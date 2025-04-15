package model

import (
	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Login struct {
	Form
}

func NewLogin() Login {
	return Login{Form: NewForm()}
}

func (l *Login) Update(msg tea.Msg) (Login, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "enter", "down":
			l.Focus = (l.Focus + 1) % 2
		case "up":
			l.Focus = (l.Focus - 1 + 2) % 2
		case "esc":
            l.reset()
            return *l, tea.Quit
		}
	}

	// Update field focus
	if l.Focus == 0 {
		l.Username.Focus()
		l.Password.Blur()

        l.Username, cmd = l.Username.Update(msg)
        l.UsernameErr = utils.ValidateUsername(l.Username.Value())
	} else {
		l.Username.Blur()
		l.Password.Focus()

        l.Password, cmd = l.Password.Update(msg)
        l.PasswordErr = utils.ValidatePassword(l.Password.Value())
	}

	return *l, cmd
}

func (m Login) View() string {
	s := "Login\n\n"
	s += "Username: " + m.Username.View() + "\n"
	if m.UsernameErr != nil {
		s += fmt.Sprintf("   [!] %s\n", m.UsernameErr.Error())
	}
	s += "Password: " + m.Password.View() + "\n"
	if m.PasswordErr != nil {
		s += fmt.Sprintf("   [!] %s\n", m.PasswordErr.Error())
	}
	s += "\n[Tab/↑↓/Enter] Switch  •  [Esc] Back"
	return s
}

func (l Login) reset() Login {
	form := NewLogin()
	return form
}
