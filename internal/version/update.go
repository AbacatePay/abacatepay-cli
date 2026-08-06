package version

import (
	"context"

	"github.com/creativeprojects/go-selfupdate"
)

// CheckUpdate checks whether a newer CLI version is available on GitHub.
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
