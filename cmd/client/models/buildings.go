package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BuildingsModel struct {
	MenuIndex int
    jar *cookiejar.Jar
}

func NewBuildingsMenu(jar *cookiejar.Jar) BuildingsModel {
    model := BuildingsModel{
        MenuIndex: 0,
        jar: jar,
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
			if m.MenuIndex < 1 {
				m.MenuIndex++
			}
		case "enter":
		}
	}

    return m, nil
}

func (m *BuildingsModel) View() string {
    costMetal, err := getCost(m.jar, "metalExtractor")
    if err != nil {
        return "Error" 
    }

    costGas, err := getCost(m.jar, "gasExtractor")
    if err != nil {
        return "Error" 
    }

    cursor := func(i int) string {
        if m.MenuIndex == i {
            return "➜ "
        }
        return "  "
    }

    cost := func(i BuildingCost) string {
        return fmt.Sprintf(
            "\tCost\n\t[Metal] :: %d\n\t[Gas] :: %d\n",
            i.MetalCost, i.GasCost,
        )
    }

    return fmt.Sprintf(
        "[Buildings]\n\n%sMetal Extractor\n%s\n%sGas Extractor\n%s\n\n",
        cursor(0), cost(*costMetal), cursor(1), cost(*costGas),
    )
}

type BuildingCost struct {
    GasCost   int64 `json:"gasCost"`
    MetalCost int64 `json:"metalCost"`
}

func getCost(jar *cookiejar.Jar, sufix string) (*BuildingCost, error) {

    url := "http://localhost:8080/buildings/" + sufix

    client := &http.Client{Timeout: 10 * time.Second, Jar: jar}

	resp, err := client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	defer resp.Body.Close()

	var cost BuildingCost

	if err := json.NewDecoder(resp.Body).Decode(&cost); err != nil {
        return nil, fmt.Errorf("Invalid json: %w", err)
	}

	if resp.StatusCode == 200 {

        return &cost, nil
	}

	return nil, fmt.Errorf("unexpected error: %s", resp.Status)
}
