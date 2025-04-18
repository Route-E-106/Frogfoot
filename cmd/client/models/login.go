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
	StateUserMenu
	StateError
)

type loginResultMsg struct {
    username string
    token string
    success bool
    err     error
}

type Login struct {
	Form
    spinner spinner.Model
    State LoginState
    userMenuModel UserMenuModel 
}

func NewLogin() Login {
    s := spinner.New()
	s.Spinner = spinner.Dot

	return Login{
        Form: NewForm(),
        spinner: s,
    }
}

func (m Login) Update(msg tea.Msg) (Login, tea.Cmd) {
	var cmd tea.Cmd
	switch m.State {
	case StateRequest:
        switch msg := msg.(type) {
        case loginResultMsg:
            if msg.success {
                m.State = StateSucceeded
                m.userMenuModel = NewUserMenu(msg.username, msg.token)
            } else {
                m.State = StateError
            }
            return m, nil
        }

        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd 
	case StateSucceeded:
        switch msg.(type) {
        case tea.KeyMsg:
            m.State = StateUserMenu
            return m, m.userMenuModel.Tick
        }

		return m, cmd
	case StateError:
        switch msg.(type) {
        case tea.KeyMsg:
            return m, utils.BackToMenuCmd()
        }

		return m, cmd
    case StateUserMenu:
		userMenuModel, cmd := m.userMenuModel.Update(msg)
		m.userMenuModel = userMenuModel
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			if m.MenuIndex > 0 {
				m.MenuIndex--
			}
        case "enter":
			if m.MenuIndex < 1 {
				m.MenuIndex++
			} else {
                m.UsernameErr = utils.ValidateUsername(m.Username.Value())
                m.PasswordErr = utils.ValidatePassword(m.Password.Value())

                if m.UsernameErr == nil && m.PasswordErr == nil {
                    m.State = StateRequest
                    return m, tea.Batch(m.spinner.Tick, attemptLogin(m))
                }
            }
		case "up":
			m.MenuIndex = (m.MenuIndex - 1 + 2) % 2
		case "esc":
            m.reset()
            return m, utils.BackToMenuCmd()
		}
	}

	if m.MenuIndex == 0 {
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

	return m, cmd
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
	case StateUserMenu:
		return m.userMenuModel.View()
	}

	return s
}

func attemptLogin(m Login) tea.Cmd {
    username := m.Username.Value()
    password := m.Password.Value()

    return func() tea.Msg {
        token, err := simulateLogin(username, password, true)
        if err != nil {
            return loginResultMsg{success: false, err: err}
        }
        return loginResultMsg{success: true, err: nil, username: username, token: token}
    }
}

func simulateLogin(username, password string, overrideHttpRequest bool) (string, error) {
    if (overrideHttpRequest) {

        if password == "fail" {
            return "", fmt.Errorf("invalid password")
        }

        return "Override", nil
    }

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
		return "custom_token", nil
	}
	return "", fmt.Errorf("unexpected error: %s", resp.Status)
}

func (l Login) reset() Login {
	form := NewLogin()
	return form
}
