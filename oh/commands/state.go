package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/openhab-cli/openhab"
)

type StateCommand struct {
	Client *openhab.Client
}

func NewStateCommand(client *openhab.Client) cli.Command {
	stateCommand := &StateCommand{Client: client}

	return cli.Command{
		Name:    "state",
		Aliases: []string{"s"},
		Usage:   "state <item name>",
		Action:  stateCommand.Action,
	}
}

func (s *StateCommand) Action(c *cli.Context) {
	if len(c.Args()) != 1 {
		PrintStringMessageAndExit("Specify item name")
	}
	item := c.Args()[0]
	if strings.Contains(item, "*") {
		itemRegexString := strings.Replace(item, "*", ".*", -1)
		itemRegexp := regexp.MustCompile(itemRegexString)

		items, err := s.Client.ListItems()
		if err != nil {
			PrintErrorMessageAndExit(err)
		}
		for _, item := range items {
			if itemRegexp.MatchString(item.Name) {
				fmt.Printf("%s %s\n", item.Name, item.State)
			}
		}
	} else {
		state, err := s.Client.GetState(c.Args()[0])
		if err != nil {
			PrintErrorMessageAndExit(err)
		}
		fmt.Printf("%s\n", state)
	}
}
