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
	err        error
}

type appState int

const (
	stateLogin appState = iota
	stateLoading
	stateDone
)

// Simulate sending credentials over HTTP
type sendCredsMsg struct{}

func sendCreds(creds userCreds) tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		// Simulated POST request
		jsonData, _ := json.Marshal(creds)

		// Replace with your actual URL if needed
		resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode != 200 {
			return tea.Quit
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
			m.statusMsg = "Login successful!"
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
				m.creds = userCreds{
					Username: m.username.Value(),
					Password: m.password.Value(),
				}
				m.state = stateLoading
				return m, tea.Batch(m.spinner.Tick, sendCreds(m.creds))
			}
			m.focus++
			if m.focus > 1 {
				m.focus = 0
			}
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
		return fmt.Sprintf(
			"Please log in:\n\nUsername: %s\nPassword: %s\n\n(Press Tab to switch, Enter to submit)",
			m.username.View(),
			m.password.View(),
		)
	case stateLoading:
		return fmt.Sprintf("\nSending credentials... %s\n(Press q to cancel)", m.spinner.View())
	case stateDone:
		return fmt.Sprintf("\n✅ %s\n\n(Press q to quit)", m.statusMsg)
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
