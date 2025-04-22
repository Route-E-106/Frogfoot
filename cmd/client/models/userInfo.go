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

type UserMenuState int

const (
	UserMenu UserMenuState = iota
    UserBuildings
)

type UserMenuModel struct {
	State     UserMenuState
	MenuIndex int
    username  string
    resources Resources
    jar *cookiejar.Jar
}

type IncomeCommand struct {
	Income int
	Timestamp time.Time
}

type ResourceHistory struct {
    ChangeAmount int
	Incomes []IncomeCommand
}

type Resources struct {
    Metal ResourceHistory
    Gas ResourceHistory
}

type tickMsg time.Time
type tickRequestMsg time.Time

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

func (m UserMenuModel) Tick() tea.Msg {
    return tickRequestMsg(time.Now())
}

func (m UserMenuModel) Update(msg tea.Msg) (UserMenuModel, tea.Cmd) {
    switch msg.(type) {
    case tickRequestMsg:
        m.updateResources()
        return m, tea.Batch(m.tick(), m.tickRequest())
    case tickMsg:
        return m, m.tick()
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
		}
	}

    return m, nil
}

func (m *UserMenuModel) View() string {
    metal := m.resources.Metal.CalculateResources()
    gas := m.resources.Gas.CalculateResources()

    s := fmt.Sprintf("[User] %s", m.username)
    s += fmt.Sprintf("\n\n[Metal] %d [Gas] %d", metal, gas)
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
func (m UserMenuModel) tickRequest() tea.Cmd {
    return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
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
    return
}
func (history ResourceHistory) CalculateResources() int {

    currentTime := time.Now()
    commands := history.Incomes
	if len(commands) == 0 {
		return history.ChangeAmount
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Timestamp.Before(commands[j].Timestamp)
	})

	total := 0.0

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

		duration := end.Sub(start).Hours()
		if duration > 0 {
			total += duration * float64(commands[i].Income)
		}
	}

    return int(math.Floor(total)) + history.ChangeAmount
}

func getResources(jar *cookiejar.Jar) (*Resources, error) {

    type rawIncome struct {
        Income          int   `json:"income"`
        ChangeTimestamp int64 `json:"change_timestamp"`
    }

    type rawResources struct {
        Gas   []rawIncome `json:"gas"`
        Metal []rawIncome `json:"metal"`
    }

    mapToResourceHistory := func (rawIncomes []rawIncome) ResourceHistory {
        var incomes []IncomeCommand
        for _, ri := range rawIncomes {
            incomes = append(incomes, IncomeCommand{
                Income:    ri.Income,
                Timestamp: time.Unix(ri.ChangeTimestamp, 0),
            })
        }
        return ResourceHistory{Incomes: incomes}
    }

    url := "http://localhost:8080/resources"

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
        resources := &Resources{
            Gas:   mapToResourceHistory(raw.Gas),
            Metal: mapToResourceHistory(raw.Metal),
        }

        return resources, nil
	}

	return nil, fmt.Errorf("unexpected error: %s", resp.Status)
}

