package update

import (
	"errors"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

var (
	overwrite bool
)

func AddToRoot(root *cobra.Command) {
	cmd := newUpdateCommand()
	root.AddCommand(cmd)
}

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update file-name",
		Short: "Updates the debt matrix",
		Long:  "Updates the debt matrix based on expenses, transactions and base state",
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

	return cmd
}

func run(_ *cobra.Command, args []string) {
	fileName := args[0]
	manager, err := sheet.LoadManager(fileName)
	if err != nil {
		log.FatalError(err)
	}

	logLoadedData(manager)

	manager.UpdateDebtors()
	if !overwrite {
		ext := path.Ext(fileName)
		fileName = strings.TrimSuffix(fileName, ext) + "-updated" + ext
	}
	err = manager.SaveAs(fileName)
	if err != nil {
		log.FatalError(err)
	}
}

func logLoadedData(manager *sheet.Manager) {
	manager.PrintData()
}
