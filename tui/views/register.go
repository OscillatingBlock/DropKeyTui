package views

import (
	"fmt"
	"strings"

	"Drop-Key-TUI/api"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type readyToRegister string

type RegisterState string

const (
	selectingMethod RegisterState = "selecting method"
	generatingKey   RegisterState = "generating key"
	enterKeyFile    RegisterState = "enter key file"
	registering     RegisterState = "registering"
)

type RegisterModel struct {
	CurrentState  RegisterState
	List          list.Model
	Inputs        []string
	statusMessage string
	err           error
	user          api.User
}

type RegistrationSuccessMsg struct {
	ID string
}

type RegistrationErrorMsg struct {
	err error
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewRegisterModel() RegisterModel {
	items := []list.Item{
		item{title: "Generate a new key pair (recommended)", desc: "Creates a new secure key pair for you."},
		item{title: "Use an existing private key file", desc: "Import an existing private key file."},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to set up your account?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))

	return RegisterModel{
		CurrentState: selectingMethod,
		List:         l,
	}
}

func (m RegisterModel) Init() tea.Cmd {
	return nil
}

func (m *RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			switch m.CurrentState {
			case selectingMethod:
				selected := m.List.SelectedItem().(item)
				if selected.title == "Generate a new key pair (recommended)" {
					m.CurrentState = generatingKey
					return m, m.generateKeyCmd()
				} else {
					m.CurrentState = enterKeyFile
					m.Inputs = []string{""}
					return m, nil
				}
			case enterKeyFile:
				m.CurrentState = registering
				return m, m.registerWithFileCmd()
			}
		case "up", "k":
			m.List.CursorUp()
			return m, nil
		case "down", "j":
			m.List.CursorDown()
			return m, nil
		}
	case api.RegisterUserResponse:
		m.CurrentState = RegisterState(done)
		m.statusMessage = fmt.Sprintf("Registration successful. User ID: %s", m.user.ID)
		return m, func() tea.Msg {
			return RegistrationSuccessMsg{ID: msg.ID}
		}
	case RegistrationErrorMsg:
		m.CurrentState = RegisterState(err)
		m.err = msg.err
		m.statusMessage = fmt.Sprintf("Registration failed: %v", m.err)
		m.CurrentState = selectingMethod
		return m, nil
	}

	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m *RegisterModel) View() string {
	var b strings.Builder

	switch m.CurrentState {
	case selectingMethod:
		b.WriteString(m.List.View())
		b.WriteString("\n(Use ↑/↓ to navigate, Enter to select, Ctrl+C to quit)")
	case generatingKey:
		b.WriteString("Generating key pair...\n")
	case enterKeyFile:
		b.WriteString("Enter path to private key file: \n")
		b.WriteString(strings.Join(m.Inputs, "\n"))
		b.WriteString("\nPress Enter to submit, Ctrl+C to quit")
	case registering:
		b.WriteString("Registering...\n")
	case RegisterState(err):
		b.WriteString(m.statusMessage + "\n\nPress Enter to retry or Ctrl+C to quit")
	case RegisterState(done):
		b.WriteString(m.statusMessage + "\n\nPress Ctrl+C to quit")
	}

	return lipgloss.NewStyle().Padding(1).Render(b.String())
}

func (m *RegisterModel) generateKeyCmd() tea.Cmd {
	return func() tea.Msg {
		// TODO call function from /crypto/keys.go
		// in /crypto/keys.go
		// 1 generate keys
		// 2 encode them to base64
		// 3 save them to .config
		// 4 return the b64 pub key
		// 5 in register.go call api to register user with the pubkey
		// if succesfull return msg to main
		// else return failure msg
	}
}
