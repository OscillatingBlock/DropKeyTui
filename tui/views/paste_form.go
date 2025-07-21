package views

import (
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"

	tea "github.com/charmbracelet/bubbletea"
)

type PasteFormModel struct {
	textarea       textarea.Model
	submit         bool
	viewport       viewport.Model
	viewportActive bool
	height         int
	width          int
}

func NewPasteFormModel() *PasteFormModel {
	physicalWidth, physicalHeight, _ := term.GetSize((os.Stdout.Fd()))
	ta := textarea.New()
	ta.Placeholder = "Type your paste here..."
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetWidth(physicalWidth - 8)
	ta.SetHeight(physicalHeight - 10)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("213"))
	vp := viewport.New(physicalWidth-18, physicalHeight-10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	return &PasteFormModel{
		textarea: ta,
		viewport: vp,
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
	return textarea.Blink
}

func (m *PasteFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			m.textarea.SetValue("")

		case "esc":
			if m.textarea.Focused() {
				m.textarea.Blur()
				return m, nil
			}
			m.textarea.Focus()
			return m, nil

		case "alt+v":
			m.viewportActive = !m.viewportActive

		case "up", "k", "down", "j", "pgup", "pgdown":
			if m.viewportActive {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
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

	out += "\n\nPress Ctrl+S to submit your paste"
	return out
}

func (m *PasteFormModel) Title() string {
	return "Create Paste"
}
