package style

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	AbacateGreen = lipgloss.Color("#a7c957")
	Yellow       = lipgloss.Color("#ffe6a7")
	Brown        = lipgloss.Color("#d4a373")
	Gray         = lipgloss.Color("#808080")

	TitleStyle = lipgloss.NewStyle().
			Foreground(AbacateGreen).
			Bold(true)

	VersionStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(AbacateGreen).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(AbacateGreen).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)
)

func AbacateTheme() *huh.Theme {
	theme := huh.ThemeBase()

	theme.Focused.Title = TitleStyle
	theme.Focused.Base = lipgloss.NewStyle().BorderForeground(AbacateGreen)
	theme.Focused.SelectedOption = lipgloss.NewStyle().Foreground(Brown)

	theme.Blurred.Title = lipgloss.NewStyle().Foreground(Gray)

	return theme
}

func Select(title string, options map[string]string) (string, error) {
	var result string
	huhOptions := make([]huh.Option[string], 0, 5)

	for label, value := range options {
		huhOptions = append(huhOptions, huh.NewOption(label, value))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(huhOptions...).
				Value(&result),
		),
	).WithTheme(AbacateTheme())

	err := form.Run()
	return result, err
}

func Container(title string) {
}

