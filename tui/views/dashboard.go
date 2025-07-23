package views

import (
	"strings"

	"Drop-Key-TUI/tui/styles"

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

var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	tab       = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

type DashboardModel struct {
	activeTab     DashboardTab
	availableTabs map[DashboardTab]DashboardTabView
	token         string
	width, height int
}

type DashboardTabView interface {
	tea.Model
	Title() string
}

func (m *DashboardModel) SetToken(token string) {
	m.token = token
}

func NewDashboardModel() *DashboardModel {
	pastes := []Paste{
		{ID: "1", title: "Login Key", Desc: "A secure login paste"},
		{ID: "2", title: "Server Config", Desc: "NGINX TLS setup"},
		{ID: "3", title: "Go Tips", Desc: "Tips for Go concurrency"},
	}

	return &DashboardModel{
		availableTabs: map[DashboardTab]DashboardTabView{
			TabCreate:     NewPasteFormModel(),
			TabYourPastes: NewPasteListModel(pastes),
			TabSearch:     NewSearchModel(),
		},
	}
}

func (m *DashboardModel) SetSize(width, height int) {
	m.height = height
	m.width = width
}

func (m *DashboardModel) Init() tea.Cmd {
	return m.availableTabs[m.activeTab].Init()
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.activeTab = (m.activeTab + 1) % tabCount
			return m, m.availableTabs[m.activeTab].Init()

		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
			return m, m.availableTabs[m.activeTab].Init()

		case "ctrl+c":
			return m, tea.Quit
		}

	case requestToken:
		return m, func() tea.Msg {
			return responseToken{token: m.token}
		}
	}
	tab := m.availableTabs[m.activeTab]
	updatedTab, cmd := tab.Update(msg)
	m.availableTabs[m.activeTab] = updatedTab.(DashboardTabView)
	return m, cmd
}

func (m *DashboardModel) View() string {
	var titles []string
	for i := DashboardTab(0); i < tabCount; i++ {
		style := tab
		if i == m.activeTab {
			style = activeTab
		}
		titles = append(titles, style.Render(m.availableTabs[i].Title()))
	}

	rawTabs := lipgloss.JoinHorizontal(lipgloss.Top, titles...)
	rawTabsWidth := lipgloss.Width(rawTabs)
	gapSize := max(0, m.width-rawTabsWidth-15)
	gap := tabGap.Render(strings.Repeat(" ", gapSize))

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, rawTabs, gap)

	var b strings.Builder
	b.WriteString(row)
	b.WriteString("\n")
	b.WriteString(m.availableTabs[m.activeTab].View())

	appStyle := styles.AppStyle.
		Height(m.height - 4).Width(m.width - 4)

	ui := appStyle.Render(
		lipgloss.Place(
			m.width-4,
			m.height-4,
			lipgloss.Left,
			lipgloss.Top,
			b.String(),
		),
	)

	return ui
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
