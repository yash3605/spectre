package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yash3605/spectre/internal/conductor"
	"github.com/yash3605/spectre/internal/models"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

var (
	colorAccent = lipgloss.Color("#BD93F9")
	colorText   = lipgloss.Color("#F8F8F2")
	colorBorder = lipgloss.Color("#6272A4")
	colorSubtle = lipgloss.Color("#44475A")

	headerStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	searchBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder)

	mainSectionStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorBorder).
				Foreground(colorText)
)

var availableModules = map[models.Tab][]string{
	models.OSINT:   {"IP lookup", "Domain", "Email", "Username", "Phone"},
	models.Infosys: {"News", "Science", "History", "GeoPolitics"},
	models.Entity:  {"Person", "Organization", "Event"},
}

type Model struct {
	ActiveTab    models.Tab
	ActiveModule int
	CurrentState models.State
	Result       models.Result
	Search       string
	Width        int
	Height       int
}

func main() {
	m := Model{
		ActiveTab:    models.OSINT,
		CurrentState: models.StateIdle,
		ActiveModule: 0,
		Search:       "",
	}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("there's an error: %v", err)
		os.Exit(1)
	}
}

func (model *Model) Init() tea.Cmd {
	return nil
}

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentModules := availableModules[model.ActiveTab]
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+q":
			return model, tea.Quit
		case "tab":
			model.ActiveTab = (model.ActiveTab + 1) % models.TabCount
			model.ActiveModule = 0
			return model, nil
		case "backspace":
			if len(model.Search) > 0 {
				model.Search = model.Search[:len(model.Search)-1]
			}
			return model, nil
		case "down":
			model.ActiveModule = (model.ActiveModule + 1) % len(currentModules)
			return model, nil
		case "up":
			if model.ActiveModule > 0 {
				model.ActiveModule -= 1
			} else {
				model.ActiveModule = len(currentModules) - 1
			}
			return model, nil
		case "enter":
			model.Result = conductor.Search(model.Search, model.ActiveTab, model.ActiveModule)
			return model, nil
		default:
			if msg.Text != "" {
				model.Search = model.Search + msg.Text
				return model, nil
			}
		}
	case tea.WindowSizeMsg:
		model.Width = msg.Width
		model.Height = msg.Height
		return model, nil
	}
	return model, nil
}

func renderTab(label string, isActive bool) string {
	if isActive {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).Render(label)
	}
	return lipgloss.NewStyle().Foreground(colorSubtle).Render(label)
}

func renderHeader(activeTab models.Tab, width int) string {
	halfWidth := width / 2

	osintTab := renderTab("OSINT", activeTab == models.OSINT)
	infoTab := renderTab("INFO", activeTab == models.Infosys)
	entityTab := renderTab("ENTITY", activeTab == models.Entity)

	left := headerStyle.Width(halfWidth).Align(lipgloss.Left).Render("SPECTRE")
	right := headerStyle.Width(halfWidth).Align(lipgloss.Right).Render(osintTab + " " + infoTab + " " + entityTab)

	return left + right
}

func renderSearchBar(search string, width int) string {
	content := "> Search: " + search + "_"
	searchText := searchBarStyle.Width(width - 4).Render(content)
	return searchText
}

func renderMainArea(width int, height int, modules []string, activeModule int, data map[string]string) string {
	var lines []string
	var resultLine []string
	var style string
	for i, module := range modules {
		if i == activeModule {
			style = lipgloss.NewStyle().Foreground(colorAccent).Render(module)
		} else {
			style = lipgloss.NewStyle().Foreground(colorText).Render(module)
		}
		lines = append(lines, style)
	}

	for key, value := range data {
		line := key + ": " + value
		resultLine = append(resultLine, line)
	}

	leftColumn := mainSectionStyle.Width((width / 4) - 2).Height(height - 4).Render("\n" + strings.Join(lines, "\n\n"))
	rightColumn := mainSectionStyle.Width((width * 3 / 4) - 2).Height(height - 4).Render("Results\n\n" + strings.Join(resultLine, "\n\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
}

func (model *Model) View() tea.View {

	currentModules := availableModules[model.ActiveTab]

	ui := renderHeader(model.ActiveTab, model.Width)
	searchBar := renderSearchBar(model.Search, model.Width)
	mainArea := renderMainArea(model.Width, model.Height, currentModules, model.ActiveModule, model.Result.Data)
	v := tea.NewView(ui + "\n" + searchBar + "\n" + mainArea)
	return v
}
