package style

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ColorPalette struct {
	Green   lipgloss.Color
	Yellow  lipgloss.Color
	Brown   lipgloss.Color
	Gray    lipgloss.Color
	White   lipgloss.Color
	SoftRed lipgloss.Color
}

var Palette = ColorPalette{
	Green:   lipgloss.Color("#a7c957"),
	Yellow:  lipgloss.Color("#ffe6a7"),
	Brown:   lipgloss.Color("#d4a373"),
	Gray:    lipgloss.Color("#808080"),
	White:   lipgloss.Color("#FFFFFF"),
	SoftRed: lipgloss.Color("#ff8787"),
}

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Palette.Green).
			Bold(true)

	VersionStyle = lipgloss.NewStyle().
			Foreground(Palette.Yellow).
			Bold(true)

	CommandStyle = lipgloss.NewStyle().
			Foreground(Palette.Green).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Palette.Gray).
			Italic(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(Palette.White).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Palette.SoftRed).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Palette.Green).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)
)

func Spinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Waiting for authorization..."

	s.Start()

	return s
}

func AbacateTheme() *huh.Theme {
	t := huh.ThemeBase()
	t.Focused.Title = t.Focused.Title.Foreground(Palette.Green).Bold(true)
	t.Focused.Base = t.Focused.Base.BorderForeground(Palette.Green)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(Palette.Brown)
	t.Blurred.Title = t.Blurred.Title.Foreground(Palette.Gray)
	return t
}

func PrintTable(headers []string, rows [][]string) {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(Palette.Green)).
		Headers(headers...).
		Rows(rows...)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().
				Foreground(Palette.Green).
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
				Foreground(Palette.Green).
				Bold(true).
				Render(name) + displayAPIKey + lipgloss.NewStyle().
				Foreground(Palette.Green).
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
				Foreground(Palette.Green).
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
		fmt.Fprintf(&sb, "%s %s\n", LabelStyle.Render(label+":"), ValueStyle.Render(value))
	}
	fmt.Println(BoxStyle.Render(sb.String()))
}

func PrintError(err string) {
	fmt.Println(BoxStyle.
		BorderForeground(Palette.SoftRed).
		Padding(0, 1).
		Render(
			ErrorStyle.Render("⚠️  Error") + "\n\n" + lipgloss.NewStyle().Foreground(Palette.White).Render(err),
		))
}

func LogWebhookReceived(event, id string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s  %s %s [%s]\n",
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(timestamp),
		lipgloss.NewStyle().Foreground(Palette.Green).Bold(true).Render("-->"),
		lipgloss.NewStyle().Bold(true).Render(event),
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(id),
	)
}

func LogWebhookForwarded(statusCode int, statusText, event string) {
	timestamp := time.Now().Format("15:04:05")
	codeColor := Palette.Green
	if statusCode < 200 || statusCode >= 300 {
		codeColor = Palette.SoftRed
	}

	codeStyle := lipgloss.NewStyle().Foreground(codeColor).Bold(true)
	bracketStyle := lipgloss.NewStyle().Foreground(Palette.Gray)

	fmt.Printf("%s  %s %s%s%s %s\n",
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(timestamp),
		lipgloss.NewStyle().Foreground(Palette.Green).Bold(true).Render("<--"),
		bracketStyle.Render("["),
		codeStyle.Render(fmt.Sprintf("%d", statusCode)),
		bracketStyle.Render("]"),
		lipgloss.NewStyle().Bold(true).Render(event),
	)
}

func LogSigningSecret(secret string) {
	fmt.Printf("%s Your webhook signing secret is %s\n",
		lipgloss.NewStyle().Foreground(Palette.Green).Bold(true).Render(">"),
		lipgloss.NewStyle().Bold(true).Render(secret),
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
