package output

import (
	"encoding/json"
	"fmt"
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

type Result struct {
	Title   string
	Fields  map[string]string
	Data    any
	Headers []string
	Rows    [][]string
}

type Formatter struct {
	format Format
	mu     sync.RWMutex
}

var defaultFormatter = &Formatter{format: FormatText}

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
	defaultFormatter.mu.RLock()
	format := defaultFormatter.format
	defaultFormatter.mu.RUnlock()

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
	defaultFormatter.mu.RLock()
	format := defaultFormatter.format
	defaultFormatter.mu.RUnlock()

	switch format {
	case FormatJSON:
		printJSONError(msg)
	case FormatTable, FormatText:
		style.PrintError(msg)
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
	_ = encoder.Encode(data)
}

func printJSONError(msg string) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(map[string]any{
		"error": msg,
	})
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
