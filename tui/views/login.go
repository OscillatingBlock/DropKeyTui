package views

import (
	"errors"
	"fmt"

	"Drop-Key-TUI/api"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type State string

const (
	readyToAuth    State = "ready to authenticate"
	authenticating State = "authenticating"
	err            State = "error"
	done           State = "done"
)

type Model struct {
	CurrentState  State
	Spinner       spinner.Model
	statusMessage string
	err           error
}

type LoginSuccessMsg struct {
	Token string
	User  api.User
}

type RegesitrationCompletedMsg struct{}

type AuthErrorMsg struct {
	err error
}

func NewLoginMethod() Model {
	s := spinner.New()
	s.Spinner = spinner.Monkey
	return Model{
		CurrentState: authenticating,
		Spinner:      s,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.authCmd())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case LoginSuccessMsg:
		m.CurrentState = done
		return m, func() tea.Msg {
			return msg
		}

	case AuthErrorMsg:
		m.CurrentState = err
		m.err = msg.err
		return m, nil
	}
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.CurrentState == authenticating {
		return "Authenticating..." + m.Spinner.View()
	}
	if m.CurrentState == err {
		return fmt.Sprint("Error during authentication, error : %w", m.err)
	}
	return fmt.Sprintf("Authentication done")
}

// TODO complete this
func (m *Model) authCmd() tea.Cmd {
	return func() tea.Msg {
		return AuthErrorMsg{err: errors.New("Authentication failed")}
	}
}
