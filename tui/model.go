package tui

import (
	"Drop-Key-TUI/api"
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

type ResizableModel interface {
	tea.Model
	SetSize(width, height int)
}

type Model struct {
	state  viewState
	width  int
	height int
	err    error
	token  string
	config *config.Config
	user   api.User

	views map[viewState]ResizableModel
}

func New() *Model {
	home := views.NewHomeModel()
	login := views.NewLoginModel()

	return &Model{
		state: homeView,
		views: map[viewState]ResizableModel{
			homeView:  home,
			loginView: login,
			// TODO Add other views
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return m.views[m.state].Init()
}

// TODO
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.views[m.state].SetSize(msg.Width, msg.Height)

	case views.RegisterSelectedMsg:
		m.state = registrationView
		return m, m.views[registrationView].Init()

	case views.LoginSelectedMsg:
		m.state = loginView
		return m, m.views[loginView].Init()

	case views.LoginSuccessMsg:
		m.token = msg.Token
		m.user = msg.User
		m.state = loggedInView
		return m, m.views[loggedInView].Init()
	}
	// Always forward the message to the currently active view
	updatedView, cmd := m.views[m.state].Update(msg)

	// Store the updated model (in case it changed internal state)
	m.views[m.state] = updatedView.(ResizableModel)

	return m, cmd
}

func (m *Model) View() string {
	return m.views[m.state].View()
}
