package views

import (
	"encoding/json"
	"fmt"
	"os"

	"Drop-Key-TUI/api"
	"Drop-Key-TUI/config"
	"Drop-Key-TUI/crypt"
	"Drop-Key-TUI/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type pasteWithTitle struct {
	api.Paste
	Title_ string
	Desc   string
}

type pasteItem struct {
	pasteWithTitle
}

func (p pasteItem) Title() string       { return p.Title_ }
func (p pasteItem) Description() string { return p.Desc }
func (p pasteItem) FilterValue() string { return p.Title_ }

type pasteListState string

const (
	showList        pasteListState = "show pastes list"
	decryptingPaste pasteListState = "decrypting paste"
	viewingPaste    pasteListState = "Viewing paste"
)

type PasteListModel struct {
	currentState   pasteListState
	list           list.Model
	pastes         []api.Paste
	spinner        spinner.Model
	viewport       viewport.Model
	selected       *pasteWithTitle
	currentPasteID string
	selectedIndex  int
	publicKey      string
}

type DecryptedPasteMsg struct {
	ID        string
	Title     string
	PlainText string
	Err       error
}

func NewPasteListModel() *PasteListModel {
	physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
	const defaultWidth = 20
	l := list.New(nil, list.NewDefaultDelegate(), defaultWidth, physicalHeight-10)
	l.Title = styles.HeaderStyle.Render("📄 Your Pastes ")
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cfg, err := config.Load()
	if err != nil {
    fmt.Println(err)
		cfg = &config.Config{}
	}
	publicKey := cfg.PublicKey

	vp := viewport.New(physicalWidth-18, physicalHeight-10)
	vp.Style = styles.VpStyle

	return &PasteListModel{
		viewport:     vp,
		spinner:      s,
		publicKey:    publicKey,
		currentState: showList,
		list:         l,
		pastes:       nil,
	}
}

func (m *PasteListModel) Init() tea.Cmd {
    cfg, err := config.Load()
    if err != nil {
        fmt.Println("Failed to load config:", err)
        return nil
    }

    m.publicKey = cfg.PublicKey
    return tea.Batch(api.GetPastes(m.publicKey), m.spinner.Tick)
}


func (m *PasteListModel) UpdateViewportContent(paste string) {
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

	str, err := renderer.Render(paste)
	if err != nil {
		m.viewport.SetContent("error while rendering glamour")
		return
	}

	m.viewport.SetContent(str)
}

func (m *PasteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.currentState == decryptingPaste {
		switch msg := msg.(type) {
		case DecryptedPasteMsg:
			if msg.Err != nil {
				m.currentState = showList
				return m, nil
			}

			m.selected = &pasteWithTitle{
				Paste:  api.Paste{ID: msg.ID},
				Title_: msg.Title,
				Desc:   msg.PlainText,
			}

			m.UpdateViewportContent(msg.PlainText)
			m.currentState = viewingPaste
			m.currentPasteID = msg.ID
			return m, nil

		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	switch m.currentState {
	case viewingPaste:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.selected = nil
				m.currentState = showList
			}

			switch msg.String() {
			case "esc":
				m.selected = nil
				m.currentState = showList
			case "up", "k":
				m.viewport.ScrollUp(1)
			case "down", "j":
				m.viewport.ScrollDown(1)
			case "pgup":
				m.viewport.ScrollUp(m.viewport.Height)
			case "pgdown":
				m.viewport.ScrollDown(m.viewport.Height)
			}
		}
		return m, nil

	case showList:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if i, ok := m.list.SelectedItem().(pasteItem); ok {
					m.currentState = decryptingPaste
					return m, decryptPasteCmd(i)
				}
      case "ctrl+r":
          cfg, err := config.Load()
          if err != nil {
              fmt.Println("Failed to reload config:", err)
              return m, nil
          }

          m.publicKey = cfg.PublicKey
          return m, api.GetPastes(m.publicKey)

        }
		case api.PasteListFetchedMsg:
			m.pastes = msg.List

			items := make([]list.Item, len(msg.List))
			for i, p := range msg.List {
				items[i] = pasteItem{
					pasteWithTitle: pasteWithTitle{
						Paste:  p,
						Title_: msg.Titles[i], // use title from the parallel slice
						Desc:   "Press Enter to view",
					},
				}
			}

			m.list.SetItems(items)
			return m, nil

		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *PasteListModel) viewSelectedPaste() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("210")).Render(m.selected.Title_)
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press Esc to go back")
	return fmt.Sprintf("📋 %s\n%s\n%s", title, m.viewport.View(), help)
}

func (m *PasteListModel) View() string {
	switch m.currentState {
	case decryptingPaste:
		text := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).
			Render("🔐 Decrypting paste...")
		return fmt.Sprintf("\n%s %s\n", m.spinner.View(), text)

	case viewingPaste:
		if m.selected != nil {
			return m.currentPasteID + "\n" + m.viewSelectedPaste()
		}
		return "⚠️ No paste selected"

	default:
		help := styles.HelpStyle.Render("j k , h l, arrow keys to navigate | Ctrl+R to refresh")
		return "\n" + m.list.View() + "\n" + help
	}
}

func (m *PasteListModel) Title() string {
	if m.selected != nil {
		return "Paste Detail"
	}
	return "Paste List"
}

func decryptPasteCmd(p pasteItem) tea.Cmd {
	return func() tea.Msg {
		plain, err := crypt.DecryptPaste(p.ID, p.Ciphertext)
		if err != nil {
			fmt.Printf("error while decrypting : %v", err)
			return DecryptedPasteMsg{Err: err}
		}

		var data struct {
			Title string `json:"title"`
			Paste string `json:"paste"`
		}
		if err := json.Unmarshal([]byte(plain), &data); err != nil {
			return DecryptedPasteMsg{
				ID:        p.ID,
				Title:     "Invalid JSON",
				PlainText: "",
				Err:       err,
			}
		}

		return DecryptedPasteMsg{
			ID:        p.ID,
			Title:     data.Title, // use decrypted title
			PlainText: data.Paste, // only paste body
			Err:       nil,
		}
	}
}
