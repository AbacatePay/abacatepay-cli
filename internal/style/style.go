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
	t := huh.ThemeBase()

	t.Focused.Title = t.Focused.Title.Foreground(AbacateGreen).Bold(true)
	t.Focused.Base = t.Focused.Base.BorderForeground(AbacateGreen)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(Brown)

	t.Blurred.Title = t.Blurred.Title.Foreground(Gray)

	return t
}

func Select(title string, options map[string]string) (string, error) {
	var result string
	huhOptions := make([]huh.Option[string], 0, len(options))

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

func Input(title, placeholder string, value *string, validate func(string) error) error {
	input := huh.NewInput().
		Title(title).
		Placeholder(placeholder).
		Value(value)

	if validate != nil {
		input.Validate(validate)
	}

	form := huh.NewForm(huh.NewGroup(input)).WithTheme(AbacateTheme())
	return form.Run()
}

func Confirm(title string, value *bool) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Value(value),
		),
	).WithTheme(AbacateTheme())

	return form.Run()
}

