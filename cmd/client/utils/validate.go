package utils

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

type BackToMenuMsg struct{}

func ValidateUsername(username string) error {
	if len(username) < 4 {
		return errors.New("Username must be at least 4 characters")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 4 {
		return errors.New("Password must be at least 4 characters")
	}
	return nil
}

func BackToMenuCmd() tea.Cmd {
	return func() tea.Msg {
		return BackToMenuMsg{}
	}
}
