package views

import tea "github.com/charmbracelet/bubbletea"

type SearchModel struct{}

func NewSearchModel() SearchModel {
	return SearchModel{}
}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SearchModel) View() string {
	return "🔍 Search Pastes Page (Mock)"
}

func (m SearchModel) Title() string {
	return "Search Pastes"
}
