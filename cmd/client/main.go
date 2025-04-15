// based on https://github.com/charmbracelet/bubbletea-app-template
package main

import (
	"bytes"
	"encoding/json"
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

type model struct {
	username   textinput.Model
	password   textinput.Model
	focus      int
	state      appState
	spinner    spinner.Model
	creds      userCreds
	statusMsg  string
	errMsg     string
	err        error
}

type appState int

const (
	stateLogin appState = iota
	stateLoading
	stateDone
	stateError
)

type sendCredsMsg struct{}
type invalidCredsMsg struct{}

func sendCreds(creds userCreds) tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		// Simulate invalid login
		if creds.Password == "fail" {
			return invalidCredsMsg{}
		}

		// Simulate HTTP POST (replace with real API if needed)
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
	u.Focus()
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
		focus:    0,
		state:    stateLogin,
		spinner:  s,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateLogin:
		return m.updateLogin(msg)
	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		switch msg := msg.(type) {
		case sendCredsMsg:
			m.state = stateDone
			m.statusMsg = "✅ Login successful!"
		case invalidCredsMsg:
			m.state = stateError
			m.errMsg = "❌ Invalid credentials. Try again."
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
		}
		return m, cmd
	case stateDone:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "q" || key.String() == "ctrl+c") {
			return m, tea.Quit
		}
	case stateError:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "enter" {
				// Reset to login state
				m.username.SetValue("")
				m.password.SetValue("")
				m.username.Focus()
				m.password.Blur()
				m.focus = 0
				m.state = stateLogin
				m.errMsg = ""
			}
		}
	}
	return m, nil
}

func (m model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			if msg.String() == "enter" && m.focus == 1 {
				username := m.username.Value()
				password := m.password.Value()

				if len(username) < 4 || len(password) < 4 {
					m.errMsg = "Username and password must be at least 4 characters."
					return m, nil
				}

				m.creds = userCreds{Username: username, Password: password}
				m.state = stateLoading
				m.errMsg = ""
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
	} else {
		m.password, cmd = m.password.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateLogin:
		msg := fmt.Sprintf(
			"Please log in:\n\nUsername: %s\nPassword: %s\n\n(Press Tab to switch, Enter to submit)",
			m.username.View(),
			m.password.View(),
		)
		if m.errMsg != "" {
			msg += fmt.Sprintf("\n\n[Error] %s", m.errMsg)
		}
		return msg
	case stateLoading:
		return fmt.Sprintf("\nSending credentials... %s\n(Press q to cancel)", m.spinner.View())
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
