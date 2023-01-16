package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const message = `"No regret for the confidence betrayed
No more hiding in shadow
Cause I won't wait for the debt to be repaid
Time has come for you..."`

func ExecuteRoot() {
	fmt.Println(message)
	fmt.Println()

	root := &cobra.Command{
		Use:   "gem",
		Short: "Spreadsheet based program for managing group expenses",
	}

	if err := root.Execute(); err != nil {
		fmt.Printf("failed to execute root command: %v\n", err)
		os.Exit(1)
	}
}
