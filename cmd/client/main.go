// based on https://github.com/charmbracelet/bubbletea-app-template
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type userCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type appState int

const (
	stateMenu appState = iota
	stateLogin
	stateRegister
	stateLoading
	stateDone
	stateError
)

type model struct {
	username     textinput.Model
	password     textinput.Model
	focus        int
	state        appState
	spinner      spinner.Model
	creds        userCreds
	statusMsg    string
	errMsg       string
	menuIndex    int
	usernameErr  error
	passwordErr  error
}

type sendCredsMsg struct{}
type invalidCredsMsg struct{}

func validateUsername(username string) error {
	if len(username) < 4 {
		return errors.New("Username must be at least 4 characters")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 4 {
		return errors.New("Password must be at least 4 characters")
	}
	return nil
}

func sendCreds(creds userCreds) tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		if creds.Password == "fail" {
			return invalidCredsMsg{}
		}
		jsonData, _ := json.Marshal(creds)
		resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode != 200 {
			return invalidCredsMsg{}
		}
		return sendCredsMsg{}
	})
}

func initialModel() model {
	u := textinput.New()
	u.Placeholder = "Username"
	u.CharLimit = 64
	u.Width = 20

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	p.CharLimit = 64
	p.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		username: u,
		password: p,
		state:    stateMenu,
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	switch m.state {
	case stateMenu:
		return m.updateMenu(msg)
	case stateLogin:
		return m.updateLogin(msg)
	case stateRegister:
		return m.updateRegister(msg)
	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		switch msg.(type) {
            case sendCredsMsg:
                m.state = stateDone
                m.statusMsg = "✅ Login successful!"
            case invalidCredsMsg:
                m.state = stateError
                m.errMsg = "❌ Invalid credentials. Try again."
		}
		return m, cmd
	case stateDone:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "q" || key.String() == "ctrl+c") {
			return m, tea.Quit
		}
	case stateError:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			m.username.SetValue("")
			m.password.SetValue("")
			m.usernameErr = nil
			m.passwordErr = nil
			m.username.Focus()
			m.password.Blur()
			m.focus = 0
			m.state = stateLogin
		}
	}
	return m, nil
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.menuIndex > 0 {
				m.menuIndex--
			}
		case "down":
			if m.menuIndex < 1 {
				m.menuIndex++
			}
		case "enter":
			m.username.SetValue("")
			m.password.SetValue("")
			m.usernameErr = nil
			m.passwordErr = nil
			m.username.Focus()
			m.focus = 0
			if m.menuIndex == 0 {
				m.state = stateLogin
			} else {
				m.state = stateRegister
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateMenu
			m.username.SetValue("")
			m.password.SetValue("")
			m.usernameErr = nil
			m.passwordErr = nil
			return m, nil
		case "tab", "shift+tab", "up", "down", "enter":
			if msg.String() == "enter" && m.focus == 1 {
				m.usernameErr = validateUsername(m.username.Value())
				m.passwordErr = validatePassword(m.password.Value())
				if m.usernameErr != nil || m.passwordErr != nil {
					return m, nil
				}
				m.creds = userCreds{
					Username: m.username.Value(),
					Password: m.password.Value(),
				}
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, sendCreds(m.creds))
			}
			m.focus = (m.focus + 1) % 2
			if m.focus == 0 {
				m.username.Focus()
				m.password.Blur()
			} else {
				m.password.Focus()
				m.username.Blur()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.focus == 0 {
		m.username, cmd = m.username.Update(msg)
		m.usernameErr = validateUsername(m.username.Value())
	} else {
		m.password, cmd = m.password.Update(msg)
		m.passwordErr = validatePassword(m.password.Value())
	}
	return m, cmd
}

func (m model) updateRegister(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.state = stateMenu
			m.username.SetValue("")
			m.password.SetValue("")
			m.usernameErr = nil
			m.passwordErr = nil
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		cursor := func(i int) string {
			if m.menuIndex == i {
				return "➜ "
			}
			return "  "
		}
		return fmt.Sprintf(
			"\nWelcome! Choose an option:\n\n%sLogin\n%sRegister\n\n(Use ↑/↓ and Enter)",
			cursor(0), cursor(1),
		)
	case stateLogin:
		out := fmt.Sprintf("Login:\n\nUsername: %s\n", m.username.View())
		if m.usernameErr != nil {
			out += fmt.Sprintf("   [!] %s\n", m.usernameErr.Error())
		}
		out += fmt.Sprintf("Password: %s\n", m.password.View())
		if m.passwordErr != nil {
			out += fmt.Sprintf("   [!] %s\n", m.passwordErr.Error())
		}
		out += "\n\n(Tab to switch fields, Enter to submit, Esc to go back)"
		return out
	case stateRegister:
		return "\n🛠️  Registration screen (coming soon!)\n\n(Press Esc to go back)"
	case stateLoading:
		return fmt.Sprintf("\nSending credentials... %s\n", m.spinner.View())
	case stateDone:
		return fmt.Sprintf("\n%s\n\n(Press q to quit)", m.statusMsg)
	case stateError:
		return fmt.Sprintf("\n[Error] %s\n\n(Press Enter to try again)", m.errMsg)
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
