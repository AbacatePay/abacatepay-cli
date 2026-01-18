package version

import (
	"context"
	"fmt"
	"os"

	"abacatepay-cli/internal/style"

	"github.com/creativeprojects/go-selfupdate"
)

func CheckUpdate(ctx context.Context, currentVersion string) (*selfupdate.Release, bool, error) {
	slug := "AbacatePay/abacatepay-cli"
	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug(slug))
	if err != nil {
		return nil, false, err
	}

	if !found || currentVersion == "" || latest.LessOrEqual(currentVersion) {
		return nil, false, nil
	}

	return latest, true, nil
}

func ShowUpdate(currentVersion string) {
	latest, found, _ := CheckUpdate(context.Background(), currentVersion)

	if !found {
		return
	}

	msg := fmt.Sprintf(
		"🥑 %s %s\n      Current: %s\n\n   To update, run:\n   %s",
		style.TitleStyle.Render("Update available:"),
		style.VersionStyle.Render(latest.Version()),
		currentVersion,
		style.CommandStyle.Render("abacatepay update"),
	)

	fmt.Fprintln(os.Stderr, style.BoxStyle.Render(msg))
}
