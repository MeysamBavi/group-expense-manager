package cmd

import (
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/cmd/create"
	"github.com/MeysamBavi/group-expense-manager/internal/cmd/message"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gem",
	Short: "Spreadsheet based program for managing group expenses",
}

func init() {
	message.AddToRoot(rootCmd)
	create.AddToRoot(rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("failed to execute root command: %v\n", err)
		os.Exit(1)
	}
}
