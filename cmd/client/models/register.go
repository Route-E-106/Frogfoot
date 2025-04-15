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

func (r *Register) Update(msg tea.Msg) (Register, tea.Cmd) {
	var cmd tea.Cmd

	if r.Focus == 0 {
		r.Username, cmd = r.Username.Update(msg)
		r.UsernameErr = utils.ValidateUsername(r.Username.Value())
	} else {
		r.Password, cmd = r.Password.Update(msg)
		r.PasswordErr = utils.ValidatePassword(r.Password.Value())
	}

	return *r, cmd
}

func (m Register) View() string {
	out := fmt.Sprintf("Register:\n\nUsername: %s\n", m.Username.View())
	if m.UsernameErr != nil {
		out += fmt.Sprintf("   [!] %s\n", m.UsernameErr.Error())
	}
	out += fmt.Sprintf("Password: %s\n", m.Password.View())
	if m.PasswordErr != nil {
		out += fmt.Sprintf("   [!] %s\n", m.PasswordErr.Error())
	}
	out += "\n\n(Tab to switch, Enter to submit, Esc to go back)"
	return out
}
