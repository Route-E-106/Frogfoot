package model

import (
    "fmt"
	"math"
	"sort"
	"time"
    "net/http/cookiejar"
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
    Minerals ResourceHistory
    Gas ResourceHistory
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

func (m UserMenuModel) Tick() tea.Msg {
    return tickMsg(time.Now())
}

func (m UserMenuModel) Update(msg tea.Msg) (UserMenuModel, tea.Cmd) {
    switch msg.(type) {
    case tickMsg:
        m.updateResources()
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
    minerals := m.resources.Minerals.CalculateResources()
    gas := m.resources.Gas.CalculateResources()

    s := fmt.Sprintf("[User] %s", m.username)
    s += fmt.Sprintf("\n\n[Minerals] %d [Gas] %d", minerals, gas)
    cursor := func(i int) string {
        if m.MenuIndex == i {
            return "➜ "
        }
        return "  "
    }
    return fmt.Sprintf(
        "\n%s\n\n%sBuildings\n%sShips\n%sLogout\n(Use ↑/↓ and Enter)",
        s, cursor(0), cursor(1), cursor(2),
    )
}

func (m UserMenuModel) tick() tea.Cmd {
    return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m *UserMenuModel) updateResources() {
    minerals, err := getResourceHistory(m.jar); 
    if err != nil {
        return; 
    }
    gas, err := getResourceHistory(m.jar); 

    if err != nil {
        return;
    }

    m.resources = Resources {
        Minerals: minerals,
        Gas: gas,
    }

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

func getResourceHistory(jar *cookiejar.Jar) (ResourceHistory, error) {
    commands := []IncomeCommand{
        {Income: 1000, Timestamp: time.Date(2025, time.January, 01, 0, 0, 0, 0, time.UTC)},
        {Income: 2000, Timestamp: time.Date(2025, time.February, 01, 0, 0, 0, 0, time.UTC)},
        {Income: 10000, Timestamp: time.Date(2025, time.March, 01, 0, 0, 0, 0, time.UTC)},
	}

    return ResourceHistory{ChangeAmount: 100, Incomes: commands}, nil
}
