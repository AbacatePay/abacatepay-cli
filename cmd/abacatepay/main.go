package main

import (
	"fmt"
	"os"

	rootcmd "github.com/AbacatePay/abacatepay-cli/cmd"
	"github.com/AbacatePay/abacatepay-cli/internal/logger"
	"github.com/AbacatePay/abacatepay-cli/internal/style"
)

func main() {
	logCfg, err := logger.DefaultConfig()
	if err != nil {
		style.PrintError(fmt.Sprintf("Failed to configure logger: %v", err))
		os.Exit(1)
	}

	if _, err := logger.Setup(logCfg); err != nil {
		style.PrintError(fmt.Sprintf("Failed to initialize logger: %v", err))
		os.Exit(1)
	}

	rootcmd.Exec()
}
