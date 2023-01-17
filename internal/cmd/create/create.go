package create

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	membersFile string
	outputFile  string
)

func AddToRoot(root *cobra.Command) {
	createCmd := newCreateCommand()
	root.AddCommand(createCmd)
}

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new spreadsheet by taking members information",
		Long: `Creates a new spreadsheet to be used. The members' names and card numbers need to be entered one by one or passed in a csv file.
Format of every row in the csv file should be "name,cardNumber"`,
		Example: "create -f list.csv",
		Run:     run,
	}

	cmd.Flags().StringVarP(
		&membersFile,
		"file",
		"f",
		"",
		"specifies the csv file containing members information",
	)

	cmd.Flags().StringVarP(
		&outputFile,
		"output",
		"o",
		"expense-manager.xlsx",
		"specifies the output file name and path",
	)

	return cmd
}

func run(_ *cobra.Command, _ []string) {
	var members []*model.Member
	if membersFile == "" {
		members = getMembersFromStdin()
	} else {
		members = getMembersFromFile(membersFile)
	}

	for _, m := range members {
		fmt.Println(*m)
	}

	if len(members) < 2 {
		log.Fatal(errors.New("number of members should be more than 1"))
	}

	manager := sheet.NewManager(members)
	err := manager.SaveAs(outputFile)
	if err != nil {
		log.Fatal(err)
	}
}

func getMembersFromFile(file string) []*model.Member {
	csvFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(csvFile)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var members []*model.Member
	for i, v := range records {
		members = append(members, &model.Member{
			ID:         model.MID(i),
			Name:       v[0],
			CardNumber: v[1],
		})
	}

	return members
}

func getMembersFromStdin() []*model.Member {
	var members []*model.Member

	fmt.Println("Enter member's name and their card number followed by a space. Press enter for next member.")
	fmt.Println("If names or card numbers contain spaces, enclose them in \"\".")
	fmt.Println("End the process by entering an empty line.")

	scanner := bufio.NewScanner(os.Stdin)
	r := regexp.MustCompile("^((\"[\\w ]+\")|(\\w+))\\s+((\"[-\\w ]+\")|([-\\w]+))$")

	var i model.MID
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		groups := r.FindStringSubmatch(line)
		if groups == nil {
			fmt.Println("You input does not match the specified format. Try again.")
			continue
		}

		name := groups[1]
		cardNumber := groups[4]

		members = append(members, &model.Member{
			ID:         i,
			Name:       strings.Trim(name, " \""),
			CardNumber: strings.Trim(cardNumber, " \""),
		})

		i++
	}

	return members
}
