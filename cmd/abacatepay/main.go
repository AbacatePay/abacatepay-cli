package main

import (
	"fmt"
	"os"

	rootcmd "github.com/AbacatePay/abacatepay-cli/cmd"
	"github.com/AbacatePay/abacatepay-cli/internal/logger"
)

func main() {
	logCfg, err := logger.DefaultConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to configure logger: %v\n", err)
		os.Exit(1)
	}

	if _, err := logger.Setup(logCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	rootcmd.Exec()
}
