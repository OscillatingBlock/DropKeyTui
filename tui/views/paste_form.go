package views

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/crypt"
	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
)

type formState string

const (
	decidingTitle   formState = "deciding title"
	writingPaste    formState = "writing paste"
	selectingExpiry formState = "selecting expiry"
	formErr         formState = "form error"
	pastecreated    formState = "paste created successfully"
)

type PasteFormModel struct {
	currentState formState
	textarea     textarea.Model
	titleBar     textarea.Model
	viewport     viewport.Model

	viewportActive  bool
	selectingExpiry bool
	pasteCreated    bool

	height     int
	width      int
	expiryDays int

	pasteID  string
	pasteUrl string
	token    string
	title    string

	err    bool
	ErrMsg string
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

	cipherTextPayload struct {
		Title string `json:"title"`
		Paste string `json:"paste"`
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
	ta.SetHeight(physicalHeight - 12)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	vp := viewport.New(physicalWidth-18, physicalHeight-10)
	vp.Style = styles.VpStyle

	titleBar := textarea.New()
	titleBar.Focus()
	titleBar.ShowLineNumbers = false
	titleBar.CharLimit = 40
	titleBar.SetWidth(physicalWidth - 28)
	titleBar.SetHeight(physicalHeight - 30)
	titleBar.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	titleBar.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))

	return &PasteFormModel{
		currentState: decidingTitle,
		textarea:     ta,
		titleBar:     titleBar,
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
	return tea.Batch(m.titleBar.Cursor.BlinkCmd(),
		func() tea.Msg {
			return requestToken{}
		},
	)
}

func (m *PasteFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.currentState == pastecreated {
			m.currentState = writingPaste
			return m, nil
		}
		if m.currentState == formErr {
			m.currentState = writingPaste
			return m, nil
		}

		switch msg.String() {
		case "enter":
			if m.currentState == decidingTitle {
				m.title = m.titleBar.Value()
				m.titleBar.SetValue("")
				m.currentState = writingPaste
				m.textarea.SetValue("")
				return m, nil
			}
		case "esc":
			if m.currentState == formErr {
				m.err = false
				return m, nil
			}

			if m.currentState == writingPaste {
				if m.textarea.Focused() {
					m.textarea.Blur()
					return m, nil
				}
				m.textarea.Focus()
				return m, nil
			}

		case "alt+v":
			m.viewportActive = !m.viewportActive
			return m, nil

		case "up", "k", "down", "j", "pgup", "pgdown":
			if m.viewportActive {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}

		case "ctrl+s":
			if m.currentState == writingPaste {
				m.textarea.Blur()
				m.currentState = selectingExpiry
				return m, nil
			}

		case "1", "2", "3", "4", "5", "6", "7":
			if m.currentState == selectingExpiry {
				m.expiryDays = int(msg.Runes[0]-'0') * 86400
				paste := m.textarea.Value()
				return m, m.CreatePaste(paste, m.title, m.token, m.expiryDays)
			}

		case "alt+c":
			if m.currentState == writingPaste {
				m.textarea.SetValue("")
				return m, nil
			}

		case "alt+n":
			if m.currentState == writingPaste {
				m.currentState = decidingTitle
				m.textarea.SetValue("")
				m.titleBar.SetValue("")
				return m, nil
			}
		}

	case api.PasteCreatedMsg:
		m.pasteUrl = msg.URL
		m.pasteID = msg.ID
		m.currentState = pastecreated

		// remap tempID -> actualID
		return m, remapTempIdCmd(msg.TempID, msg.CreatePasteResponse.ID)

	case api.ErrMsg:
		m.currentState = formErr
		m.ErrMsg = msg.Error()
		return m, nil

	case responseToken:
		m.token = msg.token
	}
	if m.currentState == decidingTitle {
		m.titleBar, cmd = m.titleBar.Update(msg)
	}
	if !m.viewportActive {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

func (m *PasteFormModel) View() string {
	out := "\n"

	switch m.currentState {

	case decidingTitle:
		_, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
		out += styles.HeaderStyle.Render("📝 Enter a title for your paste:")
		out += "\n"
		out += m.titleBar.View()
		out += styles.HelpStyle.PaddingTop(physicalHeight - 12).Render("tab to switch tabs | Ctrl+C to quit")

	case writingPaste:
		if m.viewportActive {
			m.UpdateViewportContent()
			out += m.viewport.View()
		} else {
			out += "\n"
			out += m.textarea.View()
		}
		out += "\n" + m.renderHelp()

	case selectingExpiry:
		out += styles.HeaderStyle.Render("⏳ Select expiry (1–7 days):\n")
		out += styles.HelpStyle.Render("Use number keys to choose expiry")

	case pastecreated:
		urlStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Underline(true).
			MarginTop(1)

		res := styles.SuccessHeaderStyle.Render("✔ Paste created successfully")
		id := urlStyle.Render(fmt.Sprintf("🔗 Paste ID: %v", m.pasteID))
		help := styles.HelpStyle.Render("Press any key to continue...")

		out += lipgloss.JoinVertical(lipgloss.Left, res, id, help)

	case formErr:
		err := styles.ErrStyle.Render("✘ " + m.ErrMsg)
		help := styles.HelpStyle.Render("Press any key to continue...")
		out += lipgloss.JoinVertical(lipgloss.Left, err, help)

	default:
		out += styles.ErrStyle.Render("Unknown state")
	}

	return out
}

// CreatePaste creates a paste
func (m *PasteFormModel) CreatePaste(paste, title, token string, expiresIn int) tea.Cmd {
	// instead of sending cipher text directly create a json payload
	// which will have paste title and cipher text both
	plainCipherText, err := m.getCipherTextPayload(title, paste)
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}

	// use a new temp id to encrypt each new paste
	tempID := uuid.New().String()

	// encryt the payload
	encrypted, err := crypt.EncryptPaste(tempID, plainCipherText)
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}

	// encode the payload
	encB64 := base64.StdEncoding.EncodeToString([]byte(encrypted))
	user, err := config.Load()
	if err != nil {
		return func() tea.Msg {
			return api.ErrMsg(err)
		}
	}

	privKeyBytes, _ := base64.StdEncoding.DecodeString(user.PrivateKey)
	sigBytes := ed25519.Sign(privKeyBytes, []byte(encrypted))
	sigB64 := base64.StdEncoding.EncodeToString(sigBytes)

	// Call API
	return api.CreatePaste(api.PasteRequest{
		Ciphertext: encB64,
		Signature:  sigB64,
		PublicKey:  user.PublicKey,
		ExpiresIn:  expiresIn,
	},
		token,
		tempID)
}

// remapTempIdCmd remaps the TempID to actualID given by the server

func remapTempIdCmd(tempID, actualID string) tea.Cmd {
	return func() tea.Msg {
		if err := crypt.MoveKey(tempID, actualID); err != nil {
			return api.ErrMsg(err)
		}
		return nil
	}
}

// returns a cipher text payload which have fields title and paste
func (m *PasteFormModel) getCipherTextPayload(title, paste string) ([]byte, error) {
	blob := cipherTextPayload{
		Title: title,
		Paste: paste,
	}
	return json.Marshal(blob)
}

func (m *PasteFormModel) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		MarginTop(1)

	return helpStyle.Render(
		"Ctrl+S to submit | Esc to switch mode | Alt+V preview | Alt+C clear | Alt+N new paste",
	)
}

func (m *PasteFormModel) Title() string {
	return "Create Paste"
}
