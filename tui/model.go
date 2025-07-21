package tui

import (
	"os"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/tui/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type viewState int

const (
	homeView viewState = iota
	registrationView
	loginView
	dashbordView
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
	register := views.NewRegisterModel()
	dashbord := views.NewDashboardModel()

	return &Model{
		state: homeView,
		views: map[viewState]ResizableModel{
			homeView:         home,
			registrationView: register,
			loginView:        login,
			dashbordView:     dashbord,

			// TODO Add other views
		},
	}
}

func (m *Model) Init() tea.Cmd {
	physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
	m.width = physicalWidth
	m.height = physicalHeight
	m.views[m.state].SetSize(physicalWidth, physicalHeight)
	return m.views[m.state].Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.views[m.state].SetSize(msg.Width, msg.Height)

	case views.RegisterSelectedMsg:
		m.state = registrationView
		m.views[registrationView].SetSize(m.width, m.height)
		physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
		m.views[m.state].SetSize(physicalWidth, physicalHeight)
		return m, m.views[registrationView].Init()

	case views.RegistrationSuccessMsg:
		m.state = loginView
		physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
		m.views[m.state].SetSize(physicalWidth, physicalHeight)
		config.SaveUserID(msg.ID)
		return m, m.views[loginView].Init()

	case views.LoginSelectedMsg:
		m.state = loginView
		physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
		m.views[m.state].SetSize(physicalWidth, physicalHeight)
		return m, m.views[loginView].Init()

	case views.LoginSuccessMsg:
		m.token = msg.Token
		m.user = msg.User
		m.state = dashbordView
		m.views[m.state].SetSize(m.width, m.height)
		return m, tea.Batch(
			m.views[dashbordView].Init(),
			func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			},
		)

	case views.RequestUserIDMsg:
		userID, err := config.LoadUserID()
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}
		return m, func() tea.Msg {
			return views.UserID{ID: userID}
		}

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
