package views

import tea "github.com/charmbracelet/bubbletea"

type HomeModel struct{}

func (m *HomeModel) Init() tea.Cmd {
	return func() tea.Msg {
		return HomeModel{}
	}
}

type HomeRegisterSelectedMsg struct{}

type HomeLoginSelectedMsg struct{}
