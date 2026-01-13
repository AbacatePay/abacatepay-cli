package main

import (
	"fmt"
	"os"

	"abacatepay-cli/cmd"
	"abacatepay-cli/internal/logger"
)

func main() {
	logCfg, err := logger.DefaultConfig()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao configurar logger: %v\n", err)

		// TODO: Add custom status code later
		os.Exit(1)
	}

	if _, err := logger.Setup(logCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)

		os.Exit(1)
	}

	cmd.Exec()
}
