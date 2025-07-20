package views

import tea "github.com/charmbracelet/bubbletea"

type SearchModel struct{}

func NewSearchModel() SearchModel {
	return SearchModel{}
}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	return m, nil
}

func (m SearchModel) View() string {
	return "🔍 Search Pastes Page (Mock)"
}
