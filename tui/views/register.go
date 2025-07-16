package views

import (
	"Drop-Key-TUI/api"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	user api.User
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
	return RegisterModel{}
}

func (m RegisterModel) Init() tea.Cmd {
	return nil
}

func (m *RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl +c", "q":
			return m, tea.Quit
		}
	}
}

func (m *RegisterModel) View() string {
	return ""
}
