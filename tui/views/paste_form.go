package views

import tea "github.com/charmbracelet/bubbletea"

type PasteFormModel struct {
	height int
	width  int
}

func NewPasteFormModel() PasteFormModel {
	return PasteFormModel{}
}

func (m PasteFormModel) SetSize(height, width int) {
	m.height = height
	m.width = width
}

func (m PasteFormModel) Init() tea.Cmd {
	return nil
}

func (m PasteFormModel) Update(msg tea.Msg) (PasteFormModel, tea.Cmd) {
	return m, nil
}

func (m PasteFormModel) View() string {
	return "📝 Create Paste Page (Mock)"
}
