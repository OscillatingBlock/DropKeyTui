package views

import tea "github.com/charmbracelet/bubbletea"

type PasteListModel struct{}

func NewPasteListModel() *PasteListModel {
	return &PasteListModel{}
}

func (m *PasteListModel) Init() tea.Cmd {
	return nil
}

func (m *PasteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *PasteListModel) View() string {
	return "📄 Your Pastes Page (Mock)"
}

func (m *PasteListModel) Title() string {
	return "Paste List"
}
