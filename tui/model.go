package tui

import (
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/tui/views"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	homeView viewState = iota
	registrationView
	loginView
	loggedInView
)

type Model struct {
	state  viewState
	login  views.Model
	width  int
	height int
	err    error
	home   views.HomeModel
	token  string
	config *config.Config

	// TODO add other models also
}

func New() Model {
	model := Model{
		login: views.NewLoginMethod(),
		state: homeView,
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return m.home.Init()
}

// TODO
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case tea.WindowSizeMsg:

	case views.LoginSuccessMsg:

	case views.RegisterSelectedMsg:

	case views.LoginSelectedMsg:

	}
	return nil
}
