package views

import (
	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type HomeModel struct {
	quitting bool
	choices  []string
	cursor   int
	height   int
	width    int
	err      error
	token    string
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

func (m *HomeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func NewHomeModel() *HomeModel {
	return &HomeModel{
		choices: []string{"Register", "Login"},
		cursor:  0,
	}
}

func (m *HomeModel) SetToken(token string) {
	m.token = token
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
					return LoginSelectedMsg{}
				}
			} else {
				return m, func() tea.Msg {
					return RegisterSelectedMsg{}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	}
	return m, nil
}

func (m *HomeModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	myFigure := figure.NewFigure("DropKey", "banner3-D", true)
	title := styles.HeadingStyle.Render(myFigure.String())
	loginBtn := styles.ButtonStyle.MarginLeft(1).Render("Login")
	activeLoginBtn := styles.ActiveButtonStyle.MarginLeft(1).Render("Login")
	registerBtn := styles.ButtonStyle.Render("Register")
	activeRegisterBtn := styles.ActiveButtonStyle.Render("Register")

	var buttons string
	if m.cursor == 0 {
		buttons = lipgloss.JoinVertical(lipgloss.Left, activeLoginBtn, registerBtn)
	} else {
		buttons = lipgloss.JoinVertical(lipgloss.Left, loginBtn, activeRegisterBtn)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		buttons,
	)

	appStyle := styles.AppStyle.Height(m.height - 4).Width(m.width - 2)

	ui := appStyle.Render(
		lipgloss.Place(
			m.width,
			m.height-4,
			lipgloss.Center,
			lipgloss.Center,
			content,
		),
	)
	return ui
}
