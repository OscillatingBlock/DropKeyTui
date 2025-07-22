package views

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/crypt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"

	tea "github.com/charmbracelet/bubbletea"
)

type PasteFormModel struct {
	textarea        textarea.Model
	submit          bool
	viewport        viewport.Model
	viewportActive  bool
	height          int
	width           int
	selectingExpiry bool
	expiryDays      int
	pasteID         string
	pasteUrl        string
	pasteCreated    bool
	token           string
	err             bool
	ErrMsg          string
}

type (
	pasteCreated  struct{}
	requestToken  struct{}
	responseToken struct {
		token string
	}
	pasteCreateError struct {
		err string
	}
)

func NewPasteFormModel() *PasteFormModel {
	physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
	ta := textarea.New()
	ta.Placeholder = "Type your paste here..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetWidth(physicalWidth - 18)
	ta.SetHeight(physicalHeight - 10)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	vp := viewport.New(physicalWidth-18, physicalHeight-10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F25D94")).
		PaddingRight(2)

	return &PasteFormModel{
		textarea:     ta,
		viewport:     vp,
		pasteCreated: false,
	}
}

func (m *PasteFormModel) UpdateViewportContent() {
	const glamourGutter = 2
	vpWidth := m.viewport.Width
	renderWidth := vpWidth - m.viewport.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(renderWidth),
	)
	if err != nil {
		m.viewport.SetContent("error while setting glamour renderer")
		return
	}

	str, err := renderer.Render(m.textarea.Value())
	if err != nil {
		m.viewport.SetContent("error while rendering glamour")
		return
	}

	m.viewport.SetContent(str)
}

func (m *PasteFormModel) Init() tea.Cmd {
	return tea.Batch(m.textarea.Cursor.BlinkCmd(),
		func() tea.Msg {
			return requestToken{}
		},
	)
}

func (m *PasteFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.pasteCreated {
			m.pasteCreated = false
			return m, nil
		}
		if m.err {
			m.err = false
			return m, nil
		}

		switch msg.String() {
		case "esc":
			if m.err {
				m.err = false
				return m, nil
			}

			if m.textarea.Focused() {
				m.textarea.Blur()
				return m, nil
			}
			m.textarea.Focus()
			return m, nil

		case "alt+v":
			m.viewportActive = !m.viewportActive
			return m, nil

		case "up", "k", "down", "j", "pgup", "pgdown":
			if m.viewportActive {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case "ctrl+s":
			if !m.selectingExpiry {
				m.textarea.Blur()
				m.selectingExpiry = true
				return m, nil
			}

		case "1", "2", "3", "4", "5", "6", "7":
			if m.selectingExpiry {
				m.expiryDays = int(msg.Runes[0]-'0') * 86400
				m.selectingExpiry = false
				m.submit = true
				paste := m.textarea.Value()
				return m, m.CreatePaste(paste, m.expiryDays, m.token)
			}

		case "alt+c":
			m.textarea.SetValue("")
			return m, nil
		}

	case api.CreatePasteResponse:

		m.textarea.SetValue("")
		m.pasteUrl = msg.URL
		m.pasteID = msg.ID
		m.pasteCreated = true
		return m, nil

	case api.ErrMsg:
		m.err = true
		m.ErrMsg = msg.Error()
		return m, nil

	case responseToken:
		m.token = msg.token
	}

	if !m.viewportActive {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

func (m *PasteFormModel) View() string {
	out := "\n"

	if m.viewportActive {
		m.UpdateViewportContent()
		out += m.viewport.View()
	} else {
		out += m.textarea.View()
	}

	if m.pasteCreated {
		headerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Bold(true)

		urlStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Underline(true).
			MarginTop(1)

		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true).
			MarginTop(1)

		res := headerStyle.Render("✔ Paste created successfully")
		url := urlStyle.Render(fmt.Sprintf("🔗 Paste URL: %v", m.pasteUrl))
		help := helpStyle.Render("Press any key to continue...")

		out += lipgloss.JoinVertical(lipgloss.Left, res, url, help)
		return out
	}

	if m.err {
		errStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("234")).
			Padding(0, 1).
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("124")).
			MarginTop(1)
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true).
			MarginTop(1)

		err := errStyle.Render("✘ " + m.ErrMsg)
		help := helpStyle.Render("Press any key to continue...")

		out += lipgloss.JoinVertical(lipgloss.Left, err, help)

		return out
	}

	out += "\n" + m.renderHelp()
	if m.selectingExpiry {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Render("\nChoose expiry for paste (1–7 days):")
	}

	return out
}

func (m *PasteFormModel) Title() string {
	return "Create Paste"
}

func (m *PasteFormModel) CreatePaste(paste string, expiresIn int, token string) tea.Cmd {
	encryptedPaste, err := crypt.EncryptPaste(paste)
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}
	pasteB64 := base64.StdEncoding.EncodeToString([]byte(encryptedPaste))

	user, err := config.Load()
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(user.PrivateKey)
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}
	sigBytes := ed25519.Sign(privKeyBytes, []byte(encryptedPaste))
	sigB64 := base64.StdEncoding.EncodeToString(sigBytes)
	return api.CreatePaste(api.PasteRequest{
		Ciphertext: pasteB64,
		Signature:  sigB64,
		PublicKey:  user.PublicKey,
		ExpiresIn:  expiresIn,
	},
		token,
	)
}

func (m *PasteFormModel) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginTop(1)

	return helpStyle.Render(
		"Ctrl+S to submit | Esc to switch mode | Alt+V preview | Alt+C clear",
	)
}
