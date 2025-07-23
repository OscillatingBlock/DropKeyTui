package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Paste struct {
	ID    string
	title string
	Desc  string
}

type pasteItem Paste

func (p pasteItem) Title() string       { return p.title }
func (p pasteItem) Description() string { return p.Desc }
func (p pasteItem) FilterValue() string { return p.title }

type PasteListModel struct {
	list     list.Model
	pastes   []Paste
	selected *Paste
}

func NewPasteListModel(pastes []Paste) *PasteListModel {
	items := make([]list.Item, len(pastes))
	for i, p := range pastes {
		items[i] = pasteItem(p)
	}

	const defaultWidth = 20
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, 14)
	l.Title = "📄 Your Pastes"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return &PasteListModel{
		list:   l,
		pastes: pastes,
	}
}

func (m *PasteListModel) Init() tea.Cmd {
	return nil
}

func (m *PasteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.selected != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.selected = nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(pasteItem); ok {
				p := Paste(i)
				m.selected = &p
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *PasteListModel) View() string {
	if m.selected != nil {
		return m.viewSelectedPaste()
	}
	return m.list.View()
}

func (m *PasteListModel) viewSelectedPaste() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("210")).Render(m.selected.title)
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(m.selected.Desc)
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress Esc to go back")
	return fmt.Sprintf("📋 Viewing Paste\n\n%s\n\n%s\n%s", title, desc, help)
}

func (m *PasteListModel) Title() string {
	if m.selected != nil {
		return "Paste Detail"
	}
	return "Paste List"
}
