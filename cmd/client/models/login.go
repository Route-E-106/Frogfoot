package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
    "github.com/charmbracelet/bubbles/spinner"

	tea "github.com/charmbracelet/bubbletea"
)

type LoginState int

const (
	StateInput LoginState = iota
	StateRequest
	StateSucceeded
	StateError
)

type loginResultMsg struct {
    success bool
    err     error
}

type Login struct {
	Form
    spinner spinner.Model
    State LoginState
}

func NewLogin() Login {
    s := spinner.New()
	s.Spinner = spinner.Dot

	return Login{
        Form: NewForm(),
        spinner: s,
    }
}

func (m *Login) Update(msg tea.Msg) (Login, tea.Cmd) {
	var cmd tea.Cmd
	switch m.State {
	case StateRequest:
        switch msg := msg.(type) {
        case loginResultMsg:
            if msg.success {
                m.State = StateSucceeded
            } else {
                m.State = StateError
            }
            return *m, nil
        }

        m.spinner, cmd = m.spinner.Update(msg)
        return *m, cmd 
	case StateSucceeded:
        switch msg.(type) {
        case tea.KeyMsg:
            return *m, utils.BackToMenuCmd()
        }

		return *m, cmd
	case StateError:
        switch msg.(type) {
        case tea.KeyMsg:
            return *m, utils.BackToMenuCmd()
        }

		return *m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.Focus = (m.Focus + 1) % 2
        case "enter":
            if m.Focus == 0 {
                m.Focus = (m.Focus + 1) % 2
            } else {
                m.UsernameErr = utils.ValidateUsername(m.Username.Value())
                m.PasswordErr = utils.ValidatePassword(m.Password.Value())

                if m.UsernameErr == nil && m.PasswordErr == nil {
                    m.State = StateRequest
                    return *m, tea.Batch(m.spinner.Tick, attemptLogin(m))
                }
            }
		case "up":
			m.Focus = (m.Focus - 1 + 2) % 2
		case "esc":
            m.reset()
            return *m, utils.BackToMenuCmd()
		}
	}

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

	switch m.State {
    case StateInput:
        s += "\n[Tab/↑↓/Enter] Switch  •  [Esc] Back"
	case StateRequest:
        s += fmt.Sprintf("\nSending Credentials... %s", m.spinner.View())
	case StateSucceeded:
        s += "\n✅ Login successful!"
	case StateError:
        s+= "\n❌ Invalid credentials. Try again."
	}

	return s
}

func attemptLogin(m *Login) tea.Cmd {
    username := m.Username.Value()
    password := m.Password.Value()

    return func() tea.Msg {
        _, err := simulateLogin(username, password)
        if err != nil {
            return loginResultMsg{success: false, err: err}
        }
        return loginResultMsg{success: true, err: nil}
    }
}

func simulateLogin(username, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := "https://httpbin.org/post"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if password == "fail" {
		return "", fmt.Errorf("invalid password")
	}

	if resp.StatusCode == 200 {
		return "Login successful!", nil
	}
	return "", fmt.Errorf("unexpected error: %s", resp.Status)
}

func (l Login) reset() Login {
	form := NewLogin()
	return form
}
