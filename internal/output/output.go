package output

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"abacatepay-cli/internal/style"
)

type Format string

const (
	FormatText  Format = "text"
	FormatJSON  Format = "json"
	FormatTable Format = "table"
)

type Profile struct {
	Name   string
	Token  string
	Active bool
}

type Result struct {
	Title   string
	Fields  map[string]string
	Data    any
	Headers []string
	Rows    [][]string
}

type formatter struct {
	format Format
	mu     sync.RWMutex
}

var defaultFormatter = &formatter{format: FormatText}

func SetFormat(f Format) {
	defaultFormatter.mu.Lock()
	defer defaultFormatter.mu.Unlock()
	defaultFormatter.format = f
}

func GetFormat() Format {
	defaultFormatter.mu.RLock()
	defer defaultFormatter.mu.RUnlock()
	return defaultFormatter.format
}

func ParseFormat(s string) (Format, error) {
	switch s {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "table":
		return FormatTable, nil
	default:
		return "", fmt.Errorf("invalid output format: %s (valid: text, json, table)", s)
	}
}

func Print(r Result) {
	format := GetFormat()

	switch format {
	case FormatJSON:
		printJSON(r)
	case FormatTable:
		printTable(r)
	default:
		printText(r)
	}
}

func Error(msg string) {
	format := GetFormat()

	switch format {
	case FormatJSON:
		printJSONError(msg)
	default:
		style.PrintError(msg)
	}
}

func PrintProfiles(profiles []Profile, activeProfile string) {
	format := GetFormat()

	switch format {
	case FormatJSON:
		printProfilesJSON(profiles, activeProfile)
	case FormatTable:
		printProfilesTable(profiles)
	default:
		printProfilesText(profiles, activeProfile)
	}
}

func printJSON(r Result) {
	data := r.Data
	if data == nil {
		data = map[string]any{
			"title":  r.Title,
			"fields": r.Fields,
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		slog.Debug("failed to encode JSON output", "error", err)
	}
}

func printJSONError(msg string) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]any{"error": msg}); err != nil {
		slog.Debug("failed to encode JSON error output", "error", err)
	}
}

func printText(r Result) {
	style.PrintSuccess(r.Title, r.Fields)
}

func printTable(r Result) {
	if len(r.Rows) > 0 && len(r.Headers) > 0 {
		style.PrintTable(r.Headers, r.Rows)
		return
	}

	if len(r.Fields) > 0 {
		headers := make([]string, 0, len(r.Fields))
		values := make([]string, 0, len(r.Fields))

		orderedKeys := []string{"ID", "Status", "Profile", "User", "Email", "Name", "URL", "Amount"}
		seen := make(map[string]bool)

		for _, key := range orderedKeys {
			if val, ok := r.Fields[key]; ok {
				headers = append(headers, key)
				values = append(values, val)
				seen[key] = true
			}
		}

		for key, val := range r.Fields {
			if !seen[key] {
				headers = append(headers, key)
				values = append(values, val)
			}
		}

		style.PrintTable(headers, [][]string{values})
	}
}

func printProfilesJSON(profiles []Profile, activeProfile string) {
	profileData := make([]map[string]any, 0, len(profiles))
	for _, p := range profiles {
		profileData = append(profileData, map[string]any{
			"name":   p.Name,
			"active": p.Active,
		})
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(map[string]any{
		"profiles": profileData,
		"active":   activeProfile,
	}); err != nil {
		slog.Debug("failed to encode profiles JSON output", "error", err)
	}
}

func printProfilesTable(profiles []Profile) {
	rows := make([][]string, 0, len(profiles))

	for _, p := range profiles {
		shortKey := formatShortToken(p.Token)
		activeMarker := ""
		if p.Active {
			activeMarker = "Yes"
		}
		rows = append(rows, []string{p.Name, shortKey, activeMarker})
	}

	style.PrintTable([]string{"Name", "API Key", "Active"}, rows)
}

func printProfilesText(profiles []Profile, activeProfile string) {
	profileMap := make(map[string]string, len(profiles))
	for _, p := range profiles {
		profileMap[p.Name] = p.Token
	}

	style.ProfileSimpleList(profileMap, activeProfile)
}

func formatShortToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 10 {
		return token
	}
	return token[:10] + "..."
}
