package style

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

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
)

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
		displayAPIKey := LabelStyle.Render(" (no API key)")
		if apiKey != "" {
			shortKey := apiKey
			if len(shortKey) > 10 {
				shortKey = shortKey[:10]
			}
			displayAPIKey = LabelStyle.Render(fmt.Sprintf(" (%s...)", shortKey))
		}

		if name != activeItem {
			fmt.Println(name + displayAPIKey)
			continue
		}

		output := lipgloss.NewStyle().
			Foreground(Palette.Green).
			Bold(true).
			Render(name) + displayAPIKey + lipgloss.NewStyle().
			Foreground(Palette.Green).
			Bold(true).
			Render("     🥑")
		fmt.Println(output)
	}
	fmt.Println("")
}

func SimpleList(items []string, activeItem string) {
	for _, item := range items {
		if item != activeItem {
			fmt.Println(item)
			continue
		}

		output := lipgloss.NewStyle().
			Foreground(Palette.Green).
			Bold(true).
			Render(item + "     🥑")
		fmt.Println(output)
	}
	fmt.Println("")
}

// RenderSuccess returns the success rendering used by PrintSuccess: a title
// line followed by "label: value" lines, no border/box. Shared with the
// interactive tui.RunTask so a spinner-then-result flow ends in the exact
// same text, drawn by the same code, instead of two disconnected passes.
func RenderSuccess(title string, fields map[string]string) string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("🥑 " + title))
	for label, value := range fields {
		fmt.Fprintf(&sb, "\n%s %s", LabelStyle.Render(label+":"), ValueStyle.Render(value))
	}
	return sb.String()
}

func PrintSuccess(title string, fields map[string]string) {
	fmt.Println(RenderSuccess(title, fields))
}

// PrintInfo prints a de-emphasized transient status line - e.g. instructions
// shown before a spinner starts. Keeps output consistently styled instead of
// falling back to a bare fmt.Println.
func PrintInfo(msg string) {
	FprintInfo(os.Stdout, msg)
}

// FprintInfo is PrintInfo for an arbitrary writer, for callers (like the
// non-interactive `listen` banner) that intentionally keep status text on
// stderr so piped stdout data stays clean.
func FprintInfo(w io.Writer, msg string) {
	fmt.Fprintln(w, LabelStyle.Render(msg))
}

// RenderError returns the error rendering used by PrintError, no border/box.
// See RenderSuccess.
func RenderError(err string) string {
	return ErrorStyle.Render("⚠️  Error:") + " " + lipgloss.NewStyle().Foreground(Palette.White).Render(err)
}

func PrintError(err string) {
	fmt.Println(RenderError(err))
}

func PrintVerifyError(expected, received string) {
	var sb strings.Builder
	sb.WriteString(ErrorStyle.Render("⚠️  Signature mismatch") + "\n\n")

	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Debug Analysis:") + "\n")
	fmt.Fprintf(&sb, "%s %s\n", LabelStyle.Render("Expected:"), ValueStyle.Render(expected))
	fmt.Fprintf(&sb, "%s %s\n\n", LabelStyle.Render("Received:"), ValueStyle.Render(received))

	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Common causes for mismatch:") + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(Palette.White).Render("1. Payload content differs (check for extra spaces, newlines, or formatting).") + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(Palette.White).Render("2. Wrong secret key used.") + "\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(Palette.White).Render("3. Timestamp manipulation."))

	fmt.Println(sb.String())
}

// RenderWebhookReceived renders a single received-webhook line. It's the
// single source of truth for that line's look, shared by LogWebhookReceived
// (plain sequential output) and the interactive listen dashboard.
func RenderWebhookReceived(event, id string) string {
	timestamp := time.Now().Format("15:04:05")
	return fmt.Sprintf("%s  %s %s [%s]",
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(timestamp),
		lipgloss.NewStyle().Foreground(Palette.Green).Bold(true).Render("-->"),
		lipgloss.NewStyle().Bold(true).Render(event),
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(id),
	)
}

// RenderWebhookForwarded renders a single forwarded-webhook line. See
// RenderWebhookReceived.
func RenderWebhookForwarded(statusCode int, statusText, event string) string {
	timestamp := time.Now().Format("15:04:05")
	codeColor := Palette.Green
	if statusCode < 200 || statusCode >= 300 {
		codeColor = Palette.SoftRed
	}

	codeStyle := lipgloss.NewStyle().Foreground(codeColor).Bold(true)
	bracketStyle := lipgloss.NewStyle().Foreground(Palette.Gray)

	return fmt.Sprintf("%s  %s %s%s%s %s",
		lipgloss.NewStyle().Foreground(Palette.Gray).Render(timestamp),
		lipgloss.NewStyle().Foreground(Palette.Green).Bold(true).Render("<--"),
		bracketStyle.Render("["),
		codeStyle.Render(fmt.Sprintf("%d", statusCode)),
		bracketStyle.Render("]"),
		lipgloss.NewStyle().Bold(true).Render(event),
	)
}

func LogWebhookReceived(event, id string) {
	fmt.Println(RenderWebhookReceived(event, id))
}

func LogWebhookForwarded(statusCode int, statusText, event string) {
	fmt.Println(RenderWebhookForwarded(statusCode, statusText, event))
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

// RenderJSON marshals data and returns it syntax-highlighted, matching the
// coloring PrintJSON writes to stdout. It's the single source of truth for
// that highlighting, shared with callers (like the interactive listen
// dashboard) that need the string instead of a direct print.
func RenderJSON(data any) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Debug("failed to marshal JSON for display", "error", err)
		return fmt.Sprint(data)
	}

	str := string(b)
	var result strings.Builder

	keyStyle := lipgloss.NewStyle().Foreground(Palette.Green)
	stringStyle := lipgloss.NewStyle().Foreground(Palette.Yellow)
	numberStyle := lipgloss.NewStyle().Foreground(Palette.Brown)
	boolStyle := lipgloss.NewStyle().Foreground(Palette.SoftRed)

	inString := false
	isKey := false

	for i := 0; i < len(str); i++ {
		char := str[i]

		if char == '"' {
			inString = !inString
			if inString {
				for j := i + 1; j < len(str); j++ {
					if str[j] == '"' {
						for k := j + 1; k < len(str); k++ {
							if str[k] == ' ' || str[k] == '\t' || str[k] == '\n' || str[k] == '\r' {
								continue
							}
							if str[k] == ':' {
								isKey = true
							}
							break
						}
						break
					}
				}

				if isKey {
					result.WriteString(keyStyle.Render("\""))
				} else {
					result.WriteString(stringStyle.Render("\""))
				}
			} else {
				if isKey {
					result.WriteString(keyStyle.Render("\""))
					isKey = false
				} else {
					result.WriteString(stringStyle.Render("\""))
				}
			}
			continue
		}

		if inString {
			if isKey {
				result.WriteString(keyStyle.Render(string(char)))
			} else {
				result.WriteString(stringStyle.Render(string(char)))
			}
			continue
		}

		if (char >= '0' && char <= '9') || char == '-' || char == '.' {
			result.WriteString(numberStyle.Render(string(char)))
			continue
		}

		if i+4 <= len(str) && (str[i:i+4] == "true" || str[i:i+4] == "null") {
			result.WriteString(boolStyle.Render(str[i : i+4]))
			i += 3
			continue
		}

		if i+5 <= len(str) && str[i:i+5] == "false" {
			result.WriteString(boolStyle.Render(str[i : i+5]))
			i += 4
			continue
		}

		result.WriteByte(char)
	}

	return result.String()
}

func PrintJSON(data any) {
	fmt.Println(RenderJSON(data))
}
