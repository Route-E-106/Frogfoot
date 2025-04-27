package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/Route-E-106/Frogfoot/cmd/client/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BuildingCost struct {
    GasCost   int64 `json:"gasCost"`
    MetalCost int64 `json:"metalCost"`
}


type BuildingsModel struct {
	MenuIndex int
    jar       *cookiejar.Jar
    CostMetal BuildingCost
    CostGas   BuildingCost
    resources utils.Resources
}

type requestResourcesMsg time.Time

func NewBuildingsMenu(jar *cookiejar.Jar, resources utils.Resources) BuildingsModel {
    costMetal, _ := getCost(jar, "metalExtractor")

    costGas, _ := getCost(jar, "gasExtractor")

    model := BuildingsModel{
        MenuIndex: 0,
        jar:       jar,
        CostMetal: *costMetal,
        CostGas:   *costGas,
        resources: resources,
    }

    return model
}

func (m BuildingsModel) Update(msg tea.Msg) (BuildingsModel, tea.Cmd) {
    metal, _ := m.resources.Metal.CalculateResources();
    gas, _ := m.resources.Gas.CalculateResources();

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up":
			if m.MenuIndex == 1 && metal >= m.CostMetal.MetalCost && gas >= m.CostMetal.GasCost {
				m.MenuIndex--
			} else {
                m.MenuIndex = 0
            }
		case "down":
			if m.MenuIndex == 0 && metal >= m.CostGas.MetalCost && gas >= m.CostGas.GasCost {
				m.MenuIndex++
			} else {
                m.MenuIndex = 1
            }
		case "enter":
            if m.MenuIndex == 0 && metal >= m.CostMetal.MetalCost && gas >= m.CostMetal.GasCost {
                upgradeBuilding(m.jar, "metalExtractor")
            }

            if m.MenuIndex == 1 && metal >= m.CostGas.MetalCost && gas >= m.CostGas.GasCost {
                upgradeBuilding(m.jar, "gasExtractor")
            }

            costMetal, _ := getCost(m.jar, "metalExtractor")
            costGas, _ := getCost(m.jar, "gasExtractor")

            m.CostMetal = *costMetal
            m.CostGas = *costGas
            return m, sendRequest()
		}
	}

    return m, nil
}

func sendRequest() tea.Cmd {
	return func() tea.Msg {
		return requestResourcesMsg(time.Now())
	}
}

func (m *BuildingsModel) View() string {

    metal, _ := m.resources.Metal.CalculateResources();
    gas, _   := m.resources.Gas.CalculateResources();

    getMetal := func() string {
        if metal >= m.CostMetal.MetalCost && gas >= m.CostMetal.GasCost {
            selectedStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(utils.Color))

            if m.MenuIndex == 0 {
                return selectedStyle.Render("➜ Metal Extractor")
            }
            return "  Metal Extractor"
        } else {
            selectedStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color("8"))

            return selectedStyle.Render("  Metal Extractor (Not enough resources)")
        }
    }

    getGas := func() string {
        if metal >= m.CostGas.MetalCost && gas >= m.CostGas.GasCost {
            selectedStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(utils.Color))

            if m.MenuIndex == 1 {
                return selectedStyle.Render("➜ Gas Extractor")
            }
            return "  Gas Extractor"
        } else {
            selectedStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color("8"))

            return selectedStyle.Render("  Gas Extractor (Not enough resources)")
        }
    }

    cost := func(i BuildingCost) string {
        return fmt.Sprintf(
            "    [Metal] :: %d\n    [Gas]   :: %d\n",
            i.MetalCost, i.GasCost,
        )
    }

    return fmt.Sprintf(
        "[Buildings Upgrade Cost]\n\n%s\n%s\n%s\n%s",
        getMetal(), cost(m.CostMetal), getGas(), cost(m.CostGas),
    )
}

func upgradeBuilding(jar *cookiejar.Jar, sufix string) (error) {
    url := "http://localhost:8080/buildings/" + sufix + "/upgrade"

    client := &http.Client{Timeout: 10 * time.Second, Jar: jar}

	_, err := client.Post(url, "", nil)

	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

    return nil
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
