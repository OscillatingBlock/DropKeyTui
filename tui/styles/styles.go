package styles

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(110).
			Height(70)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(2, 3, 2, 3).
			Width(AppStyle.GetWidth() - 30).
			Height(15).
			Align(lipgloss.Top).
			Align(lipgloss.Center)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	ActiveButtonStyle = ButtonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	SelectedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color("228")).
				Inherit(ButtonStyle)
)
