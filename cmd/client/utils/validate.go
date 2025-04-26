package utils

import (
	"errors"
	"math"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

const Color string = "#7AA2F6" 

func Hints() string {
	return "\n\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")). 
		Render(" (↑/↓ • [Enter] • [Esc])")
}

type IncomeCommand struct {
	Income int
	Timestamp time.Time
}

type ResourceHistory struct {
    TotalExpenses int64
	Incomes []IncomeCommand
}

type Resources struct {
    Metal ResourceHistory
    Gas ResourceHistory
}

func (history ResourceHistory) CalculateResources() (int64, int) {

    currentTime := time.Now()
    commands := history.Incomes
	if len(commands) == 0 {
		return history.TotalExpenses, 0
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Timestamp.Before(commands[j].Timestamp)
	})

	total := 0.0
    var income int

	for i := range commands {
		start := commands[i].Timestamp

		if currentTime.Before(start) {
			break
		}

		var end time.Time
		if i+1 < len(commands) {
			end = commands[i+1].Timestamp
		} else {
			end = currentTime
		}

		if end.After(currentTime) {
			end = currentTime
		}

        income = commands[i].Income
		duration := end.Sub(start).Hours()
		if duration > 0 {
			total += duration * float64(income)
		}
	}

    return int64(math.Floor(total)) + history.TotalExpenses, income
}


