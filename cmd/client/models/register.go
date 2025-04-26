package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type RegisterState int

const (
	RegisterStateInput RegisterState = iota
	RegisterStateRequest
	RegisterStateSucceeded
	RegisterStateError
)

type registerResultMsg struct {
    username string
    success bool
    err     error
}

type Register struct {
	Form
    spinner spinner.Model
    State RegisterState
    userMenuModel UserMenuModel 
}

func NewRegister() Register {
    s := spinner.New()
	s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(utils.Color))

	return Register{
        Form: NewForm(),
        spinner: s,
    }
}

func (m Register) Update(msg tea.Msg) (Register, tea.Cmd) {
	var cmd tea.Cmd
	switch m.State {
	case RegisterStateRequest:
        switch msg := msg.(type) {
        case registerResultMsg:
            if msg.success {
                m.State = RegisterStateSucceeded
            } else {
                m.State = RegisterStateError
            }

            return m, nil
        }

        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd 
	case RegisterStateSucceeded:
        switch msg.(type) {
        case tea.KeyMsg:
            return m, utils.BackToMenuCmd()
        }

		return m, cmd
	case RegisterStateError:
        switch msg.(type) {
        case tea.KeyMsg:
            return m, utils.BackToMenuCmd()
        }

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
                    m.State = RegisterStateRequest
                    return m, tea.Batch(m.spinner.Tick, attemptRegister(m))
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

func (m Register) View() string {
	s := "[Register]\n\n"
	s += "Username: " + m.Username.View() + "\n"
	if m.UsernameErr != nil {
		s += fmt.Sprintf("   [!] %s\n", m.RenderErr(m.UsernameErr))
	}
	s += "Password: " + m.Password.View() + "\n"
	if m.PasswordErr != nil {
		s += fmt.Sprintf("   [!] %s\n", m.RenderErr(m.PasswordErr))
	}

	switch m.State {
    case RegisterStateInput:
        s += utils.Hints()
	case RegisterStateRequest:
        s += fmt.Sprintf("\nSending Credentials... %s", m.spinner.View())
	case RegisterStateSucceeded:
        s += "\n✅ Register successful!"
	case RegisterStateError:
        s+= "\n❌ Invalid credentials. Try again."
    }

	return s
}

func attemptRegister(m Register) tea.Cmd {
    username := m.Username.Value()
    password := m.Password.Value()

    return func() tea.Msg {
        time.Sleep(2 * time.Second)
        err := register(username, password)
        if err != nil {
            return registerResultMsg{success: false, err: err}
        }
        return registerResultMsg{success: true, err: nil}
    }
}

func register(username, password string) (error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

    url := "http://localhost:8080/users/register"

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("unexpected error: %s", resp.Status)
}

func (l Register) reset() Register{
	form := NewRegister()
	return form
}
