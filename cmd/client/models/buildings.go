package model

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BuildingsModel struct {
	UserMenu     UserMenuModel
	MenuIndex int
}

func BuildingsMenu(userModel UserMenuModel) BuildingsModel {
    model := BuildingsModel{
        UserMenu:  userModel,
        MenuIndex: 0,
    }

    return model
}

func (m BuildingsModel) Update(msg tea.Msg) (BuildingsModel, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up":
			if m.MenuIndex > 0 {
				m.MenuIndex--
			}
		case "down":
			if m.MenuIndex < 2 {
				m.MenuIndex++
			}
		case "enter":
		}
	}

    return m, nil
}

func (m *BuildingsModel) View() string {
    metal, metalIncome := m.resources.Metal.CalculateResources()
    gas, gasIncome := m.resources.Gas.CalculateResources()

    s := fmt.Sprintf("[User] %s", m.username)
    s += fmt.Sprintf("\n\n[Metal] %d|%d [Gas] %d|%d", metal, metalIncome, gas, gasIncome)
    cursor := func(i int) string {
        if m.MenuIndex == i {
            return "➜ "
        }
        return "  "
    }
    return fmt.Sprintf(
        "\n%s\n\n%sBuildings\n%sShips\n%sLogout\n\n(Use ↑/↓ and Enter)",
        s, cursor(0), cursor(1), cursor(2),
    )
}
