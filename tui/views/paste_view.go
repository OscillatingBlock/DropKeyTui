package views

import tea "github.com/charmbracelet/bubbletea"

type ViewModel struct{}

func NewViewModel() ViewModel {
	return ViewModel{}
}

func (m ViewModel) Init() tea.Cmd {
	return nil
}

func (m ViewModel) Update(msg tea.Msg) (ViewModel, tea.Cmd) {
	return m, nil
}

func (m ViewModel) View() string {
	return "View pastes (Mock)"
}
