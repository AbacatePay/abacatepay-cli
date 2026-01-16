package style

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	AbacateGreen = lipgloss.Color("#a7c957")
	Yellow       = lipgloss.Color("#ffe6a7")
	Brown        = lipgloss.Color("#d4a373")
	Gray         = lipgloss.Color("#808080")
	White        = lipgloss.Color("#FFFFFF")

	TitleStyle = lipgloss.NewStyle().
			Foreground(AbacateGreen).
			Bold(true)

	VersionStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(AbacateGreen).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff8787")).
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

func PrintTable(headers []string, rows [][]string) {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(AbacateGreen)).
		Headers(headers...).
		Rows(rows...)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().
				Foreground(AbacateGreen).
				Bold(true).
				Align(lipgloss.Center)
		}
		return lipgloss.NewStyle().Padding(0, 1)
	})

	fmt.Println(t.Render())
}

func ProfileSimpleList(items map[string]string, activeItem string) {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		apiKey := items[name]
		displayAPIKey := ""

		displayAPIKey = LabelStyle.Render(" (no API key)")
		if apiKey != "" {
			shortKey := apiKey
			if len(shortKey) > 10 {
				shortKey = shortKey[:10]
			}

			displayAPIKey = LabelStyle.Render(fmt.Sprintf(" (%s...)", shortKey))
		}

		output := name + displayAPIKey
		if name == activeItem {
			output = lipgloss.NewStyle().
				Foreground(AbacateGreen).
				Bold(true).
				Render(name) + displayAPIKey + lipgloss.NewStyle().
				Foreground(AbacateGreen).
				Bold(true).
				Render("     🥑")
		}

		fmt.Println(output)
	}
	fmt.Println("")
}

func SimpleList(items []string, activeItem string) {
	for _, item := range items {
		output := item
		if item == activeItem {
			output = lipgloss.NewStyle().
				Foreground(AbacateGreen).
				Bold(true).
				Render(item + "     🥑")
		}

		fmt.Println(output)
	}
	fmt.Println("")
}

func PrintSuccess(title string, fields map[string]string) {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("🥑 "+title) + "\n\n")
	for label, value := range fields {
		sb.WriteString(fmt.Sprintf("%s %s\n", LabelStyle.Render(label+":"), ValueStyle.Render(value)))
	}
	fmt.Println(BoxStyle.Render(sb.String()))
}

func PrintError(err string) {
	fmt.Println(BoxStyle.Copy().
		BorderForeground(lipgloss.Color("#ff8787")).
		Padding(0, 1).
		Render(
			ErrorStyle.Render("⚠️  Error") + "\n\n" + lipgloss.NewStyle().Foreground(White).Render(err),
		))
}

func LogWebhookReceived(event, id string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s  %s %s [%s]\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(timestamp),
		lipgloss.NewStyle().Foreground(AbacateGreen).Bold(true).Render("-->"),
		lipgloss.NewStyle().Bold(true).Render(event),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(id),
	)
}

func LogWebhookForwarded(statusCode int, statusText, event string) {
	timestamp := time.Now().Format("15:04:05")
	codeColor := "#a7c957"
	if statusCode < 200 || statusCode >= 300 {
		codeColor = "#ff8787"
	}

	codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(codeColor)).Bold(true)
	bracketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	textStyle := lipgloss.NewStyle().Foreground(White)

	fmt.Printf("%s  %s %s%s %s%s %s\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Render(timestamp),
		lipgloss.NewStyle().Foreground(AbacateGreen).Bold(true).Render("<--"),
		bracketStyle.Render("["),
		codeStyle.Render(fmt.Sprintf("%d", statusCode)),
		textStyle.Render(statusText),
		bracketStyle.Render("]"),
		lipgloss.NewStyle().Bold(true).Render(event),
	)
}

func Select(title string, options map[string]string) (string, error) {
	var result string
	huhOptions := make([]huh.Option[string], 0, len(options))
	for label, value := range options {
		huhOptions = append(huhOptions, huh.NewOption(label, value))
	}
	form := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title(title).Options(huhOptions...).Value(&result))).WithTheme(AbacateTheme())
	err := form.Run()
	return result, err
}

func Input(title, placeholder string, value *string, validate func(string) error) error {
	input := huh.NewInput().Title(title + "\n").Placeholder(placeholder).Value(value)
	if validate != nil {
		input.Validate(validate)
	}
	form := huh.NewForm(huh.NewGroup(input)).WithTheme(AbacateTheme())
	return form.Run()
}

func Confirm(title string, value *bool) error {
	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(title).Value(value))).WithTheme(AbacateTheme())
	return form.Run()
}
