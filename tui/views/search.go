package views

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/crypt"
	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type SearchState string

const (
	enterID      SearchState = "enter paste id"
	viewPaste    SearchState = "view paste"
	searchErr    SearchState = "err"
	StateFetched SearchState = "paste fetched"
)

type SearchModel struct {
	state      SearchState
	ti         textinput.Model
	vp         viewport.Model
	decrypted  string
	fetched    bool
	laoding    bool
	notFound   bool
	expired    bool
	invalidKey bool
	loading    bool

	pasteID   string
	rawCipher string
	publicKey string
	signature string
	expiresAt time.Time
}

func NewSearchModel() *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Enter paste URL"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	vp := viewport.New(80, 20)
	vp.Style = styles.VpStyle
	vp.SetContent("")

	return &SearchModel{
		state: enterID,
		ti:    ti,
		vp:    vp,
	}
}

func (m *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			switch m.state {
			case enterID:
				m.pasteID = m.ti.Value()
				m.loading = true
				m.fetched = false
				m.notFound = false
				m.expired = false
				m.invalidKey = false
				m.decrypted = ""
				return m, api.GetPaste(m.pasteID)

			case StateFetched:
				m.decrypted = m.DecryptAndVerify(
					m.rawCipher,
					m.publicKey,
					m.signature,
					m.pasteID,
				)

				var pasteData struct {
					Title string `json:"title"`
					Paste string `json:"paste"`
				}

				if err := json.Unmarshal([]byte(m.decrypted), &pasteData); err != nil {
					m.vp.SetContent("[error: invalid decrypted JSON]")
				} else {
					m.UpdateViewportContent(pasteData.Paste)
					m.state = viewPaste
				}
				m.state = viewPaste
				return m, nil
			}

		case tea.KeyEsc:
			if m.state == viewPaste {
				m.state = enterID
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		if m.state == enterID {
			m.ti, cmd = m.ti.Update(msg)
			return m, cmd
		}
		if m.state == viewPaste {
			m.vp, cmd = m.vp.Update(msg)
			return m, cmd
		}

	case api.PasteFetchedMsg:
		p := msg.Paste

		m.rawCipher = p.Ciphertext
		m.publicKey = p.PublicKey
		m.signature = p.Signature
		m.expiresAt = p.ExpiresAt
		m.fetched = true
		m.loading = false
		m.state = StateFetched

		return m, nil

	case api.ErrMsg:
		m.loading = false
		err := msg.Error()
		switch {
		case strings.Contains(err, "not found"):
			m.notFound = true
		case strings.Contains(err, "expired"):
			m.expired = true
		case strings.Contains(err, "invalid"):
			m.invalidKey = true
		}
		m.state = searchErr
		return m, nil
	}

	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m *SearchModel) UpdateViewportContent(paste string) {
	const glamourGutter = 2
	vpWidth := m.vp.Width
	renderWidth := vpWidth - m.vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(renderWidth),
	)
	if err != nil {
		m.vp.SetContent("error while setting glamour renderer")
		return
	}

	str, err := renderer.Render(paste)
	if err != nil {
		m.vp.SetContent("error while rendering glamour")
		return
	}

	m.vp.SetContent(str)
}

func (m *SearchModel) View() string {
	if m.laoding {
		return styles.SpinnerStyle.Render("Decrypting paste...") + "\n"
	}

	if m.notFound {
		return styles.ErrorStyle.Render("❌ Paste not found") + "\n" + m.ti.View()
	}

	if m.expired {
		return styles.ErrorStyle.Render("⏳ Paste expired") + "\n" + m.ti.View()
	}

	if m.invalidKey {
		return styles.ErrorStyle.Render("🔑 Invalid key or corrupted data") + "\n" + m.ti.View()
	}

	switch m.state {
	case enterID:
		return "\n" + styles.HeaderStyle.Render("🔎 Search Paste by ID") + "\n\n" + m.ti.View() +
			styles.HelpStyle.PaddingTop(m.vp.Height-2).Render("Ctrl+C to quit")

	case StateFetched:
		info := fmt.Sprintf(
			"%s\n\n%s\n\n%s\n",
			styles.MetaStyle.Render("🔐 Encrypted Paste Fetched"),
			styles.SubtleStyle.Render("Expires at: "+m.expiresAt.Format(time.RFC822)),
			styles.FaintStyle.Render("Press Enter to decrypt"),
		)

		help := styles.HelpStyle.Render("Ctrl+C to quit")
		return info + "\n" + help

	case viewPaste:
		help := styles.HelpStyle.Render("esc to return back | j, k to navigate")
		return styles.HeaderStyle.Render("📄 Decrypted Paste") + "\n\n" + m.vp.View() + "\n" + help
	}

	return m.ti.View() + styles.HelpStyle.Render("Ctrl+C to quit")
}

func (m *SearchModel) Title() string {
	return "Search Pastes"
}

// DecryptAndVerify decrypts the base64 ciphertext using your AES key,
// and verifies the signature using the provided base64 public key.
func (m *SearchModel) DecryptAndVerify(ciphertextB64, pubKeyB64, sigB64, pasteID string) string {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "[decrypt error: invalid ciphertext]"
	}

	pubKey, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return "[verify error: invalid public key]"
	}

	signature, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return "[verify error: invalid signature]"
	}

	// Verify signature on pasteID + ciphertext
	if !ed25519.Verify(pubKey, ciphertext, signature) {
		return "[verify error: signature mismatch]"
	}

	// Decrypt using your AES key
	plaintext, err := crypt.DecryptPaste(m.pasteID, ciphertextB64)
	if err != nil {
		return "[decrypt error: " + err.Error() + "]"
	}

	return string(plaintext)
}
