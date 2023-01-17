package message

import (
	"fmt"
	"github.com/spf13/cobra"
)

const message = `"No regret for the confidence betrayed
No more hiding in shadow
Cause I won't wait for the debt to be repaid
Time has come for you..."`

func AddToRoot(root *cobra.Command) {
	messageCmd := newMessageCommand()
	root.AddCommand(messageCmd)
}

func newMessageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Displays a message for the debtors",
		Run:   run,
	}

	return cmd
}

func run(_ *cobra.Command, args []string) {
	fmt.Println(message)
	fmt.Println()
}
