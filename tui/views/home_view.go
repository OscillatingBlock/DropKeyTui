package views

import (
	"fmt"
	"os"

	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	quitting bool
	choices  []string
	cursor   int
	err      error
}

type (
	RegisterSelectedMsg struct{}
	LoginSelectedMsg    struct{}
	errMsg              struct{ err error }
)

func (e errMsg) Error() string { return e.err.Error() }

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func New() HomeModel {
	return HomeModel{
		choices: []string{"Register", "Login"},
		cursor:  0,
	}
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				return m, func() tea.Msg {
					return RegisterSelectedMsg{}
				}
			} else {
				return m, func() tea.Msg {
					return LoginSelectedMsg{}
				}
			}

		}
	case errMsg:
		m.err = msg
		return m, nil

	}
	return m, nil
}

func (m HomeModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	title := styles.HeadingStyle.Render("DropKey")
	loginButton := styles.ButtonStyle.Render("Login")
	activeLoginButton := styles.ActiveButtonStyle.Render("Login")
	registerButton := styles.ButtonStyle.Render("Register")
	activeRegisterButton := styles.ActiveButtonStyle.Render("Register")

	var homeView string
	if m.cursor == 0 {
		homeView = lipgloss.JoinVertical(lipgloss.Center, title, "\n", activeLoginButton, registerButton)
	} else {
		homeView = lipgloss.JoinVertical(lipgloss.Center, title, "\n", loginButton, activeRegisterButton)
	}

	return homeView
}

func main() {
	p := tea.NewProgram(&HomeModel{})
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
