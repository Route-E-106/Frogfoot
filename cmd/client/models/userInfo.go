package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UserMenuState int

const (
	UserMenu UserMenuState = iota
    UserBuildings
)

type UserMenuModel struct {
	State     UserMenuState
	MenuIndex int
    username  string
    resources utils.Resources
    jar       *cookiejar.Jar
    BuildingsModel BuildingsModel
}

type tickMsg time.Time

func NewUserMenu(username string, jar *cookiejar.Jar) UserMenuModel {
    model := UserMenuModel{
        State:     UserMenu,
        MenuIndex: 0,
        username:  username,
        jar: jar,
    }

    model.updateResources();

    return model
}
type tickRequestMsg time.Time

func (m UserMenuModel) Tick() tea.Msg {
    return tickRequestMsg(time.Now())
}

func (m UserMenuModel) Update(msg tea.Msg) (UserMenuModel, tea.Cmd) {
    switch msg.(type) {
    case requestResourcesMsg:
        m.updateResources()
        return m, nil
    case tickRequestMsg:
        m.updateResources()
        return m, tea.Batch(m.tick(), m.tickRequest())
    case tickMsg:
        return m, m.tick()
    }
    switch m.State {

    case UserBuildings:
        switch msg := msg.(type) {
            case tea.KeyMsg:
                switch msg.String() {
                case "esc":
                    m.State = UserMenu
                    return m, nil
            }
        }

		buildingsModel, cmd := m.BuildingsModel.Update(msg)
		m.BuildingsModel = buildingsModel
		return m, cmd
    }

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
            if m.MenuIndex == 0 {
                m.State = UserBuildings
            }
		}
	}

    return m, nil
}

func (m *UserMenuModel) View() string {
    metal, metalIncome := m.resources.Metal.CalculateResources()
    gas, gasIncome := m.resources.Gas.CalculateResources()

    line := lipgloss.NewStyle().
		Foreground(lipgloss.Color(utils.Color)).
		Render(strings.Repeat("─", 50))

    var incomeStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(utils.Color))

    s := fmt.Sprintf("[User] %s", m.username)
    s += fmt.Sprintf("\n\n[Metal] %d/%s [Gas] %d/%s", metal, incomeStyle.Render(strconv.Itoa(metalIncome)), gas, incomeStyle.Render(strconv.Itoa(gasIncome)))
    s += "\n" + line + "\n"


    cursor := func(i int, text string) string {
        var selectedStyle = lipgloss.NewStyle().
            Foreground(lipgloss.Color(utils.Color))

        if m.MenuIndex == i {
            return selectedStyle.Render("➜ " + text)
        }
        return "  " + text
    }

    switch m.State {
    case UserMenu:
        return fmt.Sprintf(
            "%s\n%s\n%s\n%s" + utils.Hints(),
            s, cursor(0, "Buildings"), cursor(1, "Ships"), cursor(2, "Logout"),
        )
    case UserBuildings:
        return s + "\n" + m.BuildingsModel.View() + utils.Hints()
    }

    return s
}

func (m UserMenuModel) tickRequest() tea.Cmd {
    return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
        return tickRequestMsg(t)
    })
}

func (m UserMenuModel) tick() tea.Cmd {
    return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m *UserMenuModel) updateResources() {
    resources, err := getResources(m.jar); 
    if err != nil {
        return; 
    }

    m.resources = *resources 
    m.BuildingsModel = NewBuildingsMenu(m.jar, *resources)

    return
}

func getResources(jar *cookiejar.Jar) (*utils.Resources, error) {

    type rawIncome struct {
        Income          int   `json:"income"`
        ChangeTimestamp int64 `json:"change_timestamp"`
    }

    type rawExpenses struct {
        TotalGasExpenses   int64 `json:"total_gas_expenses"`
        TotalMetalExpenses int64 `json:"total_metal_expenses"`
    }

    type rawResources struct {
        Gas   []rawIncome         `json:"gas"`
        Metal []rawIncome         `json:"metal"`
        TotalExpenses rawExpenses `json:"expenses"`
    }

    mapToResourceHistory := func (rawIncomes []rawIncome, totalExpenses int64) utils.ResourceHistory {
        var incomes []utils.IncomeCommand
        for _, ri := range rawIncomes {
            incomes = append(incomes, utils.IncomeCommand{
                Income:    ri.Income,
                Timestamp: time.Unix(ri.ChangeTimestamp, 0),
            })
        }
        return utils.ResourceHistory{Incomes: incomes, TotalExpenses: totalExpenses}
    }

    url := "http://localhost:8080/resources/history"

    client := &http.Client{Timeout: 10 * time.Second, Jar: jar}

	resp, err := client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	var raw rawResources

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
        return nil, fmt.Errorf("Invalid json: %w", err)
	}

	if resp.StatusCode == 200 {
        resources := &utils.Resources{
            Gas:   mapToResourceHistory(raw.Gas, raw.TotalExpenses.TotalGasExpenses),
            Metal: mapToResourceHistory(raw.Metal, raw.TotalExpenses.TotalMetalExpenses),
        }

        return resources, nil
	}

	return nil, fmt.Errorf("unexpected error: %s", resp.Status)
}
