package model

import (
	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Register struct {
	Form
}

func NewRegister() Register {
	return Register{Form: NewForm()}
}

func (m *Register) Update(msg tea.Msg) (Register, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "enter", "down":
			m.Focus = (m.Focus + 1) % 2
		case "up":
			m.Focus = (m.Focus - 1 + 2) % 2
		case "esc":
            m.reset()
            return *m, tea.Batch(utils.BackToMenuCmd())
		}
	}

	// Update field focus
	if m.Focus == 0 {
		m.Username.Focus()
		m.Password.Blur()

        m.Username, cmd = m.Username.Update(msg)
        m.UsernameErr = utils.ValidateUsername(m.Username.Value())
	} else {
		m.Username.Blur()
		m.Password.Focus()

        m.Password, cmd = m.Password.Update(msg)
        m.PasswordErr = utils.ValidatePassword(m.Password.Value())
	}

	return *m, cmd
}

func (m Register) View() string {
	s := "Register\n\n"
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

func (l Register) reset() Register {
	form := NewRegister()
	return form
}
