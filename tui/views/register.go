package views

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	selectingMethod State = "selecting method"
	generatingKey   State = "generating key"
	fetchingKeys    State = "fetching keys"
	enterKeyFile    State = "enter key file"
	registering     State = "registering"
	redirectToLogin State = "redirecting to login"
)

type RegisterModel struct {
	CurrentState  State
	List          list.Model
	ti            textinput.Model
	statusMessage string
	err           error
	width         int
	height        int
	user          api.User
	ID            string
}

type RegistrationSuccessMsg struct {
	ID string
}

type RegistrationErrorMsg struct {
	err error
}

type KeysGenerated struct {
	PublicKey string
}

type FetchedKeys struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type item struct {
	title, desc string
}

type waitingToRedirect struct {
	id string
}

func (m *RegisterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.List.SetSize(width-8, height-5)
	m.ti.Width = width - 10
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewRegisterModel() *RegisterModel {
	items := []list.Item{
		item{title: "Generate a new key pair (recommended)", desc: "Creates a new secure key pair for you."},
		item{title: "Use an existing private key file", desc: "Import an existing private key file."},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "How would you like to set up your account?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("218"))

	ti := textinput.New()
	ti.Placeholder = "Enter file location..."
	ti.Focus()

	return &RegisterModel{
		CurrentState: selectingMethod,
		List:         l,
		ti:           ti,
	}
}

func (m *RegisterModel) Init() tea.Cmd {
	return nil
}

func (m *RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit key handling
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}

		switch m.CurrentState {

		case selectingMethod:
			switch msg.String() {
			case "up", "k":
				m.List.CursorUp()
			case "down", "j":
				m.List.CursorDown()
			case "enter":
				selected := m.List.SelectedItem().(item)
				if selected.title == "Generate a new key pair (recommended)" {
					m.CurrentState = generatingKey
					return m, m.generateKeyCmd()
				}
				if selected.title == "Use an existing private key file" {
					m.CurrentState = enterKeyFile
				}

			}
			return m, nil // prevent falling through

		case enterKeyFile:
			if msg.String() == "enter" {
				m.CurrentState = fetchingKeys
				return m, m.LoadKeys(m.ti.Value())
			}

		case err:
			if msg.String() == "enter" {
				m.CurrentState = selectingMethod
				return m, nil
			}

		case done:
			if m.CurrentState == done {
				return m, func() tea.Msg {
					return RegistrationSuccessMsg{ID: m.ID}
				}
			}
		}
	}

	// Only update components relevant to the current state
	var cmds []tea.Cmd

	if m.CurrentState == selectingMethod {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.CurrentState == enterKeyFile {
		var cmd tea.Cmd
		m.ti, cmd = m.ti.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Handle custom messages
	switch msg := msg.(type) {
	case KeysGenerated:
		m.CurrentState = registering
		m.statusMessage = "generated PublicKey = " + msg.PublicKey
		return m, registerUser()

	case FetchedKeys:
		m.CurrentState = registering
		m.statusMessage = "Fetched Public Key successfully, Registering user"
		return m, registerWithFetchedKey(msg.PublicKey)

	case api.RegisterUserResponse:
		m.CurrentState = done
		m.statusMessage = fmt.Sprintf("Registration successful. User ID: %s", msg.ID)
		m.ID = msg.ID
		return m, nil

	case api.ErrMsg:
		m.CurrentState = err
		m.err = msg
		m.statusMessage = fmt.Sprintf("Registration failed: %v", m.err)
		return m, nil

	case RegistrationErrorMsg:
		m.CurrentState = err
		m.err = msg.err
		m.statusMessage = fmt.Sprintf("Registration failed: %v", m.err)
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *RegisterModel) View() string {
	var b strings.Builder

	m.ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))

	switch m.CurrentState {
	case selectingMethod:
		b.WriteString(m.List.View())
		b.WriteString("\n(Use ↑/↓ to navigate, Enter to select, Ctrl+C to quit)")
	case generatingKey:
		b.WriteString("Generating key pair...\n")
	case enterKeyFile:
		b.WriteString("Enter path to private key file: \n")
		b.WriteString(m.ti.View() + "\n")
		b.WriteString("\nPress Enter to submit, Ctrl+C to quit")
	case fetchingKeys:
		b.WriteString("Fetched Keys, Registering user...")
	case registering:
		b.WriteString(fmt.Sprintf("Registering... usr with PublicKey %v \n", m.statusMessage))
	case err:
		b.WriteString(m.statusMessage + "\n\nPress Enter to retry or Ctrl+C to quit")
	case done:
		b.WriteString(m.statusMessage + "\n\nPress Enter to Login or Ctrl+C to quit")
	}

	appStyle := styles.AppStyle.
		Height(m.height - 10).
		Width(m.width - 2)

	ui := appStyle.Render(
		lipgloss.Place(
			m.width,
			m.height-4,
			lipgloss.Left,
			lipgloss.Top,
			b.String(),
		),
	)

	return ui
}

func (m *RegisterModel) generateKeyCmd() tea.Cmd {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return func() tea.Msg {
			return RegistrationErrorMsg{err: err}
		}
	}

	cfg := &config.Config{
		PublicKey:  base64.StdEncoding.EncodeToString(pub),
		PrivateKey: base64.StdEncoding.EncodeToString(priv),
	}

	err = config.Save(cfg)
	if err != nil {
		return func() tea.Msg {
			return RegistrationErrorMsg{err: err}
		}
	}

	return func() tea.Msg {
		return KeysGenerated{
			PublicKey: cfg.PublicKey,
		}
	}
}

func registerUser() tea.Cmd {
	cfg, err := config.Load()
	if err != nil {
		return func() tea.Msg {
			return RegistrationErrorMsg{err: err}
		}
	}

	registrationCmd := api.RegisterUser(cfg.PublicKey)
	return registrationCmd
}

func (m *RegisterModel) LoadKeys(path string) tea.Cmd {
	rawConfigData, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return func() tea.Msg {
				return RegistrationErrorMsg{err: fs.ErrNotExist}
			}
		}
		return func() tea.Msg {
			return RegistrationErrorMsg{err: err}
		}
	}

	if len(rawConfigData) == 0 {
		slog.Error("keys json file is empty")
		return func() tea.Msg {
			return RegistrationErrorMsg{err: fmt.Errorf("empty json file")}
		}
	}

	var fetchedKeys FetchedKeys
	if err := json.Unmarshal(rawConfigData, &fetchedKeys); err != nil {
		slog.Error("error while decoding keys JSON%w", "error", err)
		return func() tea.Msg {
			return RegistrationErrorMsg{err: err}
		}
	}

	return func() tea.Msg {
		return fetchedKeys
	}
}

func registerWithFetchedKey(pubKey string) tea.Cmd {
	return api.RegisterUser(pubKey)
}
