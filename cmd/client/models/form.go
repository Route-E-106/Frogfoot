package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	Username    textinput.Model
	Password    textinput.Model
	MenuIndex       int
	UsernameErr error
	PasswordErr error
}

func NewForm() Form {
	u := textinput.New()
	u.Placeholder = "Username"
	u.CharLimit = 64
	u.Width = 20
	u.Focus()

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.CharLimit = 64
	p.Width = 20

	return Form{
		Username: u,
		Password: p,
		MenuIndex:    0,
	}
}

func (f *Form) RenderErr(error error) (string) {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	return style.Render(error.Error())
}

func (f *Form) SwitchFocus() {
	f.MenuIndex = (f.MenuIndex + 1) % 2
	if f.MenuIndex == 0 {
		f.Username.Focus()
		f.Password.Blur()
	} else {
		f.Username.Blur()
		f.Password.Focus()
	}
}
