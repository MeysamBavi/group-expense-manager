package update

import (
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

var (
	overwrite bool
	shortLog  bool
	longLog   bool
)

func AddToRoot(root *cobra.Command) {
	cmd := newUpdateCommand()
	root.AddCommand(cmd)
}

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update file-name",
		Short: "Updates the debt matrix",
		Long:  "Updates the debt matrix and settlements based on expenses, transactions and base state",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("no arguments passed as file name")
			}
			return nil
		},
		Run: run,
	}

	cmd.Flags().BoolVarP(
		&overwrite,
		"overwrite",
		"r",
		false,
		"if set, overwrites the existing file instead of creating a new copy",
	)

	cmd.Flags().BoolVarP(
		&shortLog,
		"short-log",
		"s",
		false,
		"logs the loaded data from spreadsheet in a short format (does not log expenses and transactions)",
	)

	cmd.Flags().BoolVarP(
		&longLog,
		"long-log",
		"l",
		false,
		"logs loaded data in a long format including expenses and transactions",
	)

	return cmd
}

func run(_ *cobra.Command, args []string) {
	fileName := args[0]
	manager, err := sheet.LoadManager(fileName)
	if err != nil {
		log.FatalError(err)
	}

	manager.UpdateDebtors()
	if shortLog {
		manager.PrintData(true)
	} else if longLog {
		manager.PrintData(false)
	}
	if !overwrite {
		ext := path.Ext(fileName)
		fileName = strings.TrimSuffix(fileName, ext) + "-updated" + ext
	}
	err = manager.SaveAs(fileName)
	if err != nil {
		log.FatalError(err)
	}

	fmt.Printf("Updated debt matrix and saved to %s\n", fileName)
}
