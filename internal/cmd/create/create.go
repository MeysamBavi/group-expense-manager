package create

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/MeysamBavi/group-expense-manager/internal/log"
	"github.com/MeysamBavi/group-expense-manager/internal/model"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/store"
	"github.com/MeysamBavi/group-expense-manager/internal/sheet/style"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var (
	membersFile string
	outputFile  string
	theme       string
)

var (
	validThemes = map[string]*style.Theme{
		"blue":   style.BlueTheme(),
		"green":  style.GreenTheme(),
		"red":    style.RedTheme(),
		"yellow": style.YellowTheme(),
		"purple": style.PurpleTheme(),
	}
	defaultTheme = "blue"
)

func AddToRoot(root *cobra.Command) {
	createCmd := newCreateCommand()
	root.AddCommand(createCmd)
}

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new spreadsheet by taking members information",
		Long: `Creates a new spreadsheet. The members' names and card numbers need to be passed in a csv file or entered one by one after a prompt.
Format of every row in the csv file should be <name>,<cardNumber>`,
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

	cmd.Flags().StringVarP(
		&theme,
		"theme",
		"t",
		defaultTheme,
		"specifies the color theme of the spreadsheet. valid values are "+strings.Join(getValidThemes(), ", "),
	)

	return cmd
}

func run(_ *cobra.Command, _ []string) {
	var members *store.MemberStore
	if membersFile == "" {
		members = getMembersFromStdin()
	} else {
		members = getMembersFromFile(membersFile)
	}

	fmt.Println("Members:")
	members.Range(func(_ int, member *model.Member) {
		fmt.Println(*member)
	})

	if members.Count() < 2 {
		log.FatalError(errors.New("number of members should be more than 1"))
	}

	manager := sheet.NewManager(members, getTheme())
	err := manager.SaveAs(outputFile)
	if err != nil {
		log.FatalError(err)
	}

	fmt.Println("Spreadsheet created successfully")
}

func getMembersFromFile(file string) *store.MemberStore {
	csvFile, err := os.Open(file)
	if err != nil {
		log.FatalError(err)
	}

	reader := csv.NewReader(csvFile)

	records, err := reader.ReadAll()
	if err != nil {
		log.FatalError(err)
	}

	members := store.NewMemberStore()
	for _, v := range records {
		err := members.AddMember(&model.Member{
			Name:       v[0],
			CardNumber: v[1],
		})
		if err != nil {
			log.FatalError(err)
		}
	}

	return members
}

func getMembersFromStdin() *store.MemberStore {
	members := store.NewMemberStore()

	fmt.Println("Enter member's name and their card number, separated by a space. Press enter for the next member.")
	fmt.Println("If names or card numbers contain spaces, enclose them in \"\".")
	fmt.Println("End the process by entering an empty line.")

	scanner := bufio.NewScanner(os.Stdin)
	lineExp := regexp.MustCompile("^((\"[\\w ]+\")|(\\w+))\\s+((\"[-\\w ]+\")|([-\\w]+))$")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		groups := lineExp.FindStringSubmatch(line)
		if groups == nil {
			fmt.Println("Your input does not match the specified format. Try again.")
			continue
		}

		name := groups[1]
		cardNumber := groups[4]

		err := members.AddMember(&model.Member{
			Name:       strings.Trim(name, " \""),
			CardNumber: strings.Trim(cardNumber, " \""),
		})
		if err != nil {
			log.FatalError(err)
		}
	}

	return members
}

func getValidThemes() []string {
	themes := make([]string, 0, len(validThemes))
	for k := range validThemes {
		themes = append(themes, k)
	}
	return themes
}

func getTheme() *style.Theme {
	theme, ok := validThemes[theme]
	if !ok {
		return validThemes[defaultTheme]
	}

	return theme
}
