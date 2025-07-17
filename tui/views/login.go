package views

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"

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
	width         int
	height        int
	err           error
}

type LoginSuccessMsg struct {
	Token string
	User  api.User
}

type AuthErrorMsg struct {
	err error
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func NewLoginModel() *Model {
	s := spinner.New()
	s.Spinner = spinner.Monkey
	return &Model{
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

	case api.ErrMsg:
		m.CurrentState = err
		m.err = msg
		return m, nil

	case api.AuthResponse:
		m.CurrentState = done
		config, err := config.Load()
		if err != nil {
			return m, func() tea.Msg {
				return AuthErrorMsg{err: err}
			}
		}
		successMsg := LoginSuccessMsg{
			Token: msg.Token,
			// TODO  send id field also in User
			User: api.User{
				PublicKey: config.PublicKey,
			},
		}

		return m, func() tea.Msg {
			return successMsg
		}

	case error:
		m.CurrentState = err
		m.err = msg
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

func (m *Model) authCmd() tea.Cmd {
	config, err := config.Load()
	if err != nil {
		return func() tea.Msg {
			return AuthErrorMsg{err: err}
		}
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		return func() tea.Msg {
			return AuthErrorMsg{err: fmt.Errorf("could not decode private key: %w", err)}
		}
	}

	challengeString := time.Now().UTC().Format(time.RFC3339)
	challengeB64 := base64.StdEncoding.EncodeToString([]byte(challengeString))

	signatureBytes := ed25519.Sign(privKeyBytes, []byte(challengeString))
	signatureB64 := base64.StdEncoding.EncodeToString(signatureBytes)

	authCmd := api.AuthenticateUser(api.AuthRequest{
		PublicKey: config.PublicKey,
		Signature: signatureB64,
		Challenge: challengeB64,
	})

	return authCmd
}
