package cmd

import (
	"fmt"

	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/style"

	"github.com/spf13/cobra"
)

var (
	logsLimit      int
	logsTypeFilter string
)

var logsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List historical webhook events from local log file",
	Long:  "Display webhook transactions recorded locally during listen sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listLogs()
	},
}

func init() {
	logsListCmd.Flags().IntVarP(&logsLimit, "limit", "n", 50, "Number of log entries to display")
	logsListCmd.Flags().StringVarP(&logsTypeFilter, "type", "t", "", "Filter by log type (webhook_received, webhook_forwarded, webhook_forward_failed, webhook_forward_error)")

	logsCmd.AddCommand(logsListCmd)
}

func listLogs() error {
	opts := logger.ReadOptions{
		Limit:      logsLimit,
		TypeFilter: logsTypeFilter,
	}

	entries, err := logger.ReadTransactionLogs(opts)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		logPath, _ := logger.GetLogFilePath()
		fmt.Printf("No transaction logs found.\n")
		fmt.Printf("Hint: Run 'abacatepay listen' to start recording webhook events.\n")
		fmt.Printf("Log file location: %s\n", logPath)
		return nil
	}

	if output.GetFormat() == output.FormatJSON {
		return printLogsJSON(entries)
	}
	return printLogsTable(entries)
}

func printLogsJSON(entries []logger.LogEntry) error {
	style.PrintJSON(map[string]any{
		"logs":  entries,
		"count": len(entries),
	})
	return nil
}

func printLogsTable(entries []logger.LogEntry) error {
	var rows [][]string

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		status := entry.Msg
		if entry.StatusCode > 0 {
			status = fmt.Sprintf("%s [%d]", entry.Msg, entry.StatusCode)
		}

		timestamp := entry.Time
		if entry.Timestamp != "" {
			timestamp = entry.Timestamp
		}

		rows = append(rows, []string{
			timestamp,
			status,
			entry.ID,
			entry.URL,
		})
	}

	style.PrintTable([]string{"Timestamp", "Type", "ID", "URL"}, rows)

	return nil
}
