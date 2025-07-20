package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardTab int

const (
	TabCreate DashboardTab = iota
	TabYourPastes
	TabSearch
	tabCount
)

type DashboardModel struct {
	currentTab    DashboardTab
	createPaste   PasteFormModel
	yourPastes    PasteListModel
	searchPastes  SearchModel
	width, height int
}

func NewDashboardModel() DashboardModel {
	return DashboardModel{
		currentTab:   TabCreate,
		createPaste:  NewPasteFormModel(),
		yourPastes:   NewPasteListModel(),
		searchPastes: NewSearchModel(),
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return nil
}

func (m DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.currentTab = (m.currentTab + 1) % tabCount
			return m, nil
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Forward to active tab
	switch m.currentTab {
	case TabCreate:
		var cmd tea.Cmd
		m.createPaste, cmd = m.createPaste.Update(msg)
		return m, cmd
	case TabYourPastes:
		var cmd tea.Cmd
		m.yourPastes, cmd = m.yourPastes.Update(msg)
		return m, cmd
	case TabSearch:
		var cmd tea.Cmd
		m.searchPastes, cmd = m.searchPastes.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m DashboardModel) View() string {
	var b strings.Builder

	tabBar := lipgloss.JoinHorizontal(lipgloss.Top,
		renderTab("Create", m.currentTab == TabCreate),
		renderTab("Your Pastes", m.currentTab == TabYourPastes),
		renderTab("Search", m.currentTab == TabSearch),
	)

	b.WriteString(tabBar + "\n\n")

	// Render current tab's view
	switch m.currentTab {
	case TabCreate:
		b.WriteString(m.createPaste.View())
	case TabYourPastes:
		b.WriteString(m.yourPastes.View())
	case TabSearch:
		b.WriteString(m.searchPastes.View())
	}

	return b.String()
}

func renderTab(name string, active bool) string {
	style := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	if active {
		style = style.Bold(true).Foreground(lipgloss.Color("205")) // shiny pink
	}

	return style.Render(name)
}
